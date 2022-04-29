package domain

import (
	"evol"
)

//CreateOrderCmd sent to Order Aggregate to  make a new order
type CreateOrderCmd struct {
	OrderId string
	BuyerId string
	Price   float32
	Goods   []string
}

func (o *CreateOrderCmd) Name() evol.CommandName {
	return "CreateOrderCmd"
}

func (o *CreateOrderCmd) TargetAggregateType() evol.AggregateType {
	return OrderAggregateType
}

func (o *CreateOrderCmd) TargetIdentity() string {
	return o.OrderId
}

type CancelOrderCmd struct {
	OrderId string
	Reason  string
}

func (c *CancelOrderCmd) Name() evol.CommandName {
	return "CancelOrderCmd"
}

func (c *CancelOrderCmd) TargetAggregateType() evol.AggregateType {
	return OrderAggregateType
}

func (c *CancelOrderCmd) TargetIdentity() string {
	return c.OrderId
}

type MakeReservationCmd struct {
	OrderId   string
	ProductId string //target aggregate identity
	Count     int
}

func (m *MakeReservationCmd) Name() evol.CommandName {
	return "MakeReservationCmd"
}

func (m *MakeReservationCmd) TargetAggregateType() evol.AggregateType {
	return StockAggregateType
}

func (m *MakeReservationCmd) TargetIdentity() string {
	return m.ProductId
}

type RollBackReservationCmd struct {
	OrderId   string
	ProductId string //target aggregate identity
	Count     int
}

func (r *RollBackReservationCmd) Name() evol.CommandName {
	return "RollBackReservationCmd"
}

func (r *RollBackReservationCmd) TargetAggregateType() evol.AggregateType {
	return StockAggregateType
}

func (r *RollBackReservationCmd) TargetIdentity() string {
	return r.ProductId
}

type PayOrderCmd struct {
	PaymentId string
	OrderId   string
	BuyerId   string
	Amount    float32
}

func (p *PayOrderCmd) Name() evol.CommandName {
	return "PayOrderCmd"
}

func (p *PayOrderCmd) TargetAggregateType() evol.AggregateType {
	return PaymentAggregateType
}

func (p *PayOrderCmd) TargetIdentity() string {
	return p.PaymentId
}

type OrderConfirmedCmd struct {
	OrderId string
}

func (o *OrderConfirmedCmd) Name() evol.CommandName {
	return "OrderConfirmedCmd"
}

func (o *OrderConfirmedCmd) TargetAggregateType() evol.AggregateType {
	return OrderAggregateType
}

func (o *OrderConfirmedCmd) TargetIdentity() string {
	return o.OrderId
}

func init() {
	evol.RegisterCommand(&CreateOrderCmd{})
	evol.RegisterCommand(&CancelOrderCmd{})
	evol.RegisterCommand(&MakeReservationCmd{})
	evol.RegisterCommand(&RollBackReservationCmd{})
	evol.RegisterCommand(&PayOrderCmd{})
	evol.RegisterCommand(&OrderConfirmedCmd{})
}
