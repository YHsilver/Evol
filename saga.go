package evol

import "context"

//SagaHandler is a special aggregate handle saga codec and send command
type SagaHandler interface {
	SagaType() string
	// SagaIdentity The property in the codec that will provide the value to find the Saga instance.
	//Typically, this value is an aggregate identifier of an aggregate that a specific saga monitors.
	SagaIdentity() string
	// HandleSagaEvent handle codec and send command
	HandleSagaEvent(ctx context.Context, event Event, bus CommandHandler) error

	StartSaga()
	EndSaga()
	IsAlive() bool
}

type SagaRepo interface {
	Load(identity string) SagaHandler
	Save(handler SagaHandler) error
	Delete(identity string) error
}
