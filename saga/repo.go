package saga

import (
	"evol"
	"sync"
)

type memorySagaRepo struct {
	sagas   map[string]evol.SagaHandler
	sagasMu sync.RWMutex
}

func NewMemorySagaRepo() *memorySagaRepo {
	return &memorySagaRepo{
		sagas: make(map[string]evol.SagaHandler),
	}
}

func (r *memorySagaRepo) Load(identity string) evol.SagaHandler {
	r.sagasMu.RLock()
	defer r.sagasMu.RUnlock()

	return r.sagas[identity]
}

func (r *memorySagaRepo) Save(handler evol.SagaHandler) error {
	r.sagasMu.Lock()
	defer r.sagasMu.Unlock()

	r.sagas[handler.SagaIdentity()] = handler
	return nil
}

func (r *memorySagaRepo) Delete(identity string) error {
	r.sagasMu.Lock()
	defer r.sagasMu.Unlock()
	delete(r.sagas, identity)
	return nil
}
