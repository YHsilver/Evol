package memory

import (
	"context"
	"errors"
	"evol"
	"sync"
)

type ModelRepo struct {
	db   map[string]evol.Entity
	dbMu sync.RWMutex
}

func (m *ModelRepo) Find(ctx context.Context, id string) (evol.Entity, error) {
	m.dbMu.RLock()
	defer m.dbMu.RUnlock()

	e, ok := m.db[id]
	if !ok {
		return nil, evol.ErrAggregateNotFound
	}
	return e, nil
}

func (m *ModelRepo) FindAll(ctx context.Context) ([]evol.Entity, error) {
	m.dbMu.RLock()
	defer m.dbMu.RUnlock()

	res := make([]evol.Entity, 0)
	for _, entity := range m.db {
		res = append(res, entity)
	}
	return res, nil

}

func (m *ModelRepo) Save(ctx context.Context, entity evol.Entity) error {
	m.dbMu.Lock()
	defer m.dbMu.Unlock()

	id := entity.EntityIdentity()
	if id == "" {
		return errors.New("missing entity identity")
	}
	m.db[id] = entity

	return nil

}

func (m *ModelRepo) Remove(ctx context.Context, id string) error {
	m.dbMu.Lock()
	defer m.dbMu.Unlock()

	if _, ok := m.db[id]; ok {
		delete(m.db, id)
	}

	return errors.New("no such entity")
}
