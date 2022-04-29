package evol

import (
	"context"
)

type EventCodec interface {
	MarshalEvent(context.Context, Event) ([]byte, error)

	UnmarshalEvent(context.Context, []byte) (Event, error)
}

type CommandCodec interface {
	MarshalCommand(context.Context, Command) ([]byte, error)

	UnmarshalCommand(context.Context, []byte) (Command, error)
}
