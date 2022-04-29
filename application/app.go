package application

import (
	"context"
	"evol"
	"evol/command"

	"evol/saga"
)

func Run(ctx context.Context,
	cmdBus evol.CommandBus,
	eventBus evol.EventBus,
	AggregateStore evol.AggregateStore,
	SagaStore evol.SagaRepo) error {
	//0. init dependency
	evol.CmdBus = cmdBus

	//1. register aggregate command handler
	//程序启动时才指定cmdbus， 需要将所有的aggregateType保存起来，程序bootstrap的时候，注册aggregate cmd handler
	err := RegisterCmdHandler(cmdBus, AggregateStore, eventBus)
	if err != nil {
		return err
	}
	//2. Register Event Handler
	//2.1 aggregatestore codec handlers
	for topic, handlers := range evol.EventHandlers {
		for _, handler := range handlers {
			eventBus.RegisterHandler(ctx, topic, handler)
		}
	}
	//2.2 saga
	err = saga.PrepareSagas(ctx, eventBus, cmdBus, SagaStore)
	if err != nil {
		return err
	}

	return nil
}

func RegisterCmdHandler(cmdBus evol.CommandBus, store evol.AggregateStore, evtBus evol.EventBus) error {
	cmds := evol.GetAllCmds()
	for name, cmd := range cmds {
		cmdHandler, err := command.NewAggCmdHandler(cmd.TargetAggregateType(), store, evtBus)
		if err != nil {
			return err
		}
		err = cmdBus.RegisterCmdHandler(name, cmdHandler)
		if err != nil {
			return err
		}
	}
	return nil
}
