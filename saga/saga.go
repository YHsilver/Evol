package saga

import (
	"context"
	"errors"
	"evol"
	"fmt"
	"github.com/thoas/go-funk"
	"sync"
)

type BaseSaga struct {
	sagaId   string
	isAlive  bool
	sagaType string
}

func (s *BaseSaga) SagaType() string {
	return s.sagaType
}

func (s *BaseSaga) SagaIdentity() string {
	return s.sagaId
}

func (s *BaseSaga) StartSaga() {
	s.isAlive = true
}

func (s *BaseSaga) EndSaga() {
	s.isAlive = false
}

func (s *BaseSaga) IsAlive() bool {
	return s.isAlive
}

func NewBaseSaga(id string, sagaType string) *BaseSaga {
	return &BaseSaga{
		sagaId:   id,
		isAlive:  false,
		sagaType: sagaType,
	}

}

// Sagas stores all saga manager, prepare for dependency injection when startup
// dependency: saga aggregatestore, command CmdBus, codec handler register on codec CmdBus
var sagaManagers = make(map[string]*SagaManager)
var sagaMsmu sync.RWMutex

//SagaManager manage a type of saga, means that a single SagaManager can only manage a single SagaType, implement EventHandler
type SagaManager struct {
	SagaType    string
	SagaRepo    evol.SagaRepo
	SagaFactory func(sagaType string, sagaIdentity string) evol.SagaHandler

	// Resolver decode unique sagaId from codec on all events
	// ensure the value is the same in all saga events in one saga lifecycle
	Resolver func(e evol.Event) string

	StartEvents []evol.Topic
	OnEvents    []evol.Topic
	EndEvents   []evol.Topic
	CmdBus      evol.CommandHandler
}

type SagaManagerOpt struct {
	SagaType string

	StartEvents []evol.Topic
	OnEvents    []evol.Topic
	EndEvents   []evol.Topic

	Resolver    func(e evol.Event) string
	SagaFactory func(sagaType string, sagaIdentity string) evol.SagaHandler
}

func NewSagaManager(option *SagaManagerOpt) *SagaManager {
	m := &SagaManager{
		SagaType:    option.SagaType,
		SagaFactory: option.SagaFactory,
		Resolver:    option.Resolver,
		StartEvents: option.StartEvents,
		OnEvents:    option.OnEvents,
		EndEvents:   option.EndEvents,
	}
	return m
}

// HandleEvent handle saga codec
func (m *SagaManager) HandleEvent(ctx context.Context, e evol.Event) error {
	sagaId := m.Resolver(e)

	saga := m.SagaRepo.Load(sagaId)
	//creat if not exist, start saga
	if saga == nil {
		saga = m.SagaFactory(m.SagaType, sagaId)
		saga.StartSaga()
	}

	if !saga.IsAlive() {
		return fmt.Errorf("saga: %v not alive", saga)
	}
	err := saga.HandleSagaEvent(ctx, e, m.CmdBus)

	//end saga for some specific events
	if funk.Contains(m.EndEvents, e.Topic()) {
		m.endSaga(saga)
	}

	if err != nil {
		return err
	}
	return nil
}

func (m *SagaManager) endSaga(saga evol.SagaHandler) {
	saga.EndSaga()
}

func checkParam(m *SagaManager) error {
	if m.SagaType == "" {
		return errors.New("SagaType missing")
	}

	if m.SagaFactory == nil {
		return errors.New("SagaFactory missing")
	}

	if m.StartEvents == nil || len(m.StartEvents) == 0 {
		return errors.New("StartEvents missing")
	}

	if m.EndEvents == nil || len(m.EndEvents) == 0 {
		return errors.New("EndEvents missing")
	}

	if m.Resolver == nil {
		return errors.New("codec identity Resolver  missing")
	}

	return nil
}

func RegisterSaga(m *SagaManager) error {
	if err := checkParam(m); err != nil {
		return err
	}

	sagaMsmu.Lock()
	defer sagaMsmu.Unlock()
	sagaType := m.SagaType
	if _, ok := sagaManagers[sagaType]; ok {
		return fmt.Errorf("saga %s already register", sagaType)
	}
	sagaManagers[sagaType] = m
	return nil
}

func PrepareSagas(ctx context.Context, evtBus evol.EventBus, cmdBus evol.CommandHandler, repo evol.SagaRepo) error {
	if repo == nil {
		return errors.New("missing saga aggregatestore")
	}
	sagaMsmu.RLock()
	defer sagaMsmu.RUnlock()
	for _, saga := range sagaManagers {
		saga.CmdBus = cmdBus
		saga.SagaRepo = repo
		events := append(append(saga.StartEvents, saga.OnEvents...), saga.EndEvents...)
		for _, topic := range events {
			err := evtBus.RegisterHandler(ctx, topic, saga)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type SagaFuncHandler func(ctx context.Context, event evol.Event, bus evol.CommandHandler) error

func (s SagaFuncHandler) IsAlive() bool {
	return true
}

func (s SagaFuncHandler) StartSaga() {

}

func (s SagaFuncHandler) EndSaga() {

}

func (s SagaFuncHandler) SagaType() string {
	return "saga_func_handler"
}

func (s SagaFuncHandler) SagaIdentity() string {
	return "saga_func_handler"
}

func (s SagaFuncHandler) HandleSagaEvent(ctx context.Context, event evol.Event, bus evol.CommandHandler) error {
	return s(ctx, event, bus)
}
