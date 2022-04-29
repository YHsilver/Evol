package evol

import (
	"context"
	"errors"
	"sync"
)

var CmdBus CommandBus

var cmds = make(map[CommandName]Command)
var cmdsMu sync.RWMutex

func RegisterCommand(cmd Command) error {
	cmdsMu.Lock()
	defer cmdsMu.Unlock()

	if _, ok := cmds[cmd.Name()]; ok {
		return errors.New("command already registered")
	}

	cmds[cmd.Name()] = cmd
	return nil
}

func SendCommand(ctx context.Context, cmd Command) error {
	err := CmdBus.HandleCommand(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func NewCommand(name CommandName) Command {
	cmdsMu.RLock()
	defer cmdsMu.RUnlock()

	return cmds[name]
}

func GetAllCmds() map[CommandName]Command {
	cmdsMu.RLock()
	defer cmdsMu.RUnlock()

	mapCopy := map[CommandName]Command{}

	for index, element := range cmds {
		mapCopy[index] = element
	}

	return mapCopy
}
