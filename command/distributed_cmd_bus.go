package command

import (
	"bytes"
	"context"
	"evol"
	"net/http"
)

//TODO: send command to another microservice
type Service struct {
	Url         string
	ContentType string
}

type CommandBus struct {
	Codec  evol.CommandCodec
	Router map[evol.CommandName]Service
}

func (c *CommandBus) HandleCommand(ctx context.Context, command evol.Command) error {
	service := c.Router[command.Name()]
	cmdData, err := c.Codec.MarshalCommand(ctx, command)
	if err != nil {
		return err
	}
	_, err = http.Post(service.Url, service.ContentType, bytes.NewBuffer(cmdData))
	return err
}

func (c *CommandBus) RegisterCmdHandler(cmd evol.CommandName, handler evol.CommandHandler) error {
	return nil
}
