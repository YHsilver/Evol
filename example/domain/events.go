package domain

import (
	"evol"
	"time"
)

const (
	OrderCreatedEventTopic evol.Topic = "OrderCreatedEvent"

	ProductReservedEventTopic      evol.Topic = "ProductReservedEvent"
	ProductReserveFailedEventTopic evol.Topic = "ProductReserveFailedEvent"

	OrderPayedEventTopic     evol.Topic = "OrderPayedEvent"
	OrderPayFailedEventTopic evol.Topic = "OrderPayFailedEvent"

	OrderConfirmedEventTopic evol.Topic = "OrderConfirmedEvent"
	OrderCanceledEventTopic  evol.Topic = "OrderCanceledEvent"
)

type OrderCreatedEvent struct {
	OrderId    string
	BuyerId    string
	TotalPrice float32
	ProductIds []string
}

type ProductReservedEvent struct {
	ProductId string
	OrderId   string
	Count     int
}

type ProductReserveFailedEvent struct {
	ProductId string
	OrderId   string
	Count     int
	Reason    string
}

type OrderPayedEvent struct {
	PaymentId string
	OrderId   string
	Amount    float32
	Time      time.Time
}

type OrderPayFailedEvent struct {
	PaymentId string
	OrderId   string
	Reason    string
	Time      time.Time
}
