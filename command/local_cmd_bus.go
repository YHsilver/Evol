package command

import (
	"context"
	"errors"
	"evol"
	"sync"
)

//LocalCommandBus dispatch command to aggregate, also a type of  evol.CommandHandler
//TODO: distributed command bus
type LocalCommandBus struct {
	handlers   map[evol.CommandName]evol.CommandHandler
	handlersMu sync.RWMutex
}

func NewCommandBus() *LocalCommandBus {
	return &LocalCommandBus{
		handlers: make(map[evol.CommandName]evol.CommandHandler),
	}
}

func (b *LocalCommandBus) RegisterCmdHandler(cmd evol.CommandName, handler evol.CommandHandler) error {
	b.handlersMu.Lock()
	defer b.handlersMu.Unlock()

	if _, ok := b.handlers[cmd]; ok {
		return errors.New("[evol] RegisterCmdHandler: command already registered")
	}

	b.handlers[cmd] = handler

	return nil
}

func (b *LocalCommandBus) HandleCommand(ctx context.Context, cmd evol.Command) error {
	b.handlersMu.RLock()
	defer b.handlersMu.RUnlock()

	if handler, ok := b.handlers[cmd.Name()]; ok {
		//Async command handle
		go handler.HandleCommand(ctx, cmd)
	} else {
		return errors.New("[evol] LocalCommandBus HandleCommand: command handler not found")
	}

	return nil
}
