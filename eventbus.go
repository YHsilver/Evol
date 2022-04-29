package evol

import "context"

type EventBus interface {
	EventHandler

	RegisterHandler(context.Context, Topic, EventHandler) error
}
