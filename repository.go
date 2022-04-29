package evol

import (
	"context"
	"errors"
)

var ErrAggregateNotFound = errors.New("[evol] could not find entity form repository")

//ReadRepository is a read entity repository for CQRS query
type ReadRepository interface {
	// Find returns an entity for an Identity
	Find(context.Context, string) (Entity, error)

	// FindAll returns all entities in the repository
	FindAll(context.Context) ([]Entity, error)
}

//WriteRepository is a write entity repository for CQRS command
type WriteRepository interface {
	// Save an Entity in the repository
	Save(context.Context, Entity) error

	// Remove an Entity with identity
	Remove(context.Context, string) error
}

type RWRepository interface {
	ReadRepository
	WriteRepository
}

//EventRepo is an aggregate codec repository for Event Sourcing
type EventRepo interface {
	// Save appends all events in the codec stream to the store
	Save(ctx context.Context, events []Event) error

	// Load loads all events for the aggregate id from the store
	Load(context.Context, string) ([]Event, error)
}
