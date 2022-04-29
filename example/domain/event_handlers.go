package domain

import (
	"context"
	"evol"
)

func init() {
	evol.RegisterEventHandlers(
		evol.EventHandlerFunc(OrderEventFuncHandler),
		OrderCreatedEventTopic,
		OrderPayedEventTopic,
		OrderConfirmedEventTopic,
		OrderCanceledEventTopic)
}

func OrderEventFuncHandler(ctx context.Context, event evol.Event) error {
	switch event.Topic() {
	case OrderCreatedEventTopic:
		//TODO:
	case OrderPayedEventTopic:
		//TODO:
	case OrderConfirmedEventTopic:
		//TODO:
	case OrderCanceledEventTopic:
		//TODO:
	}
	return nil
}
