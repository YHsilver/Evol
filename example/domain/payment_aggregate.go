package domain

import (
	"context"
	"evol"
	"github.com/mitchellh/mapstructure"
	"time"
)

func init() {
	// Register Aggregate
	evol.RegisterAggregate(PaymentAggregateType, func(s string) evol.Aggregate {
		return &PaymentAggregate{
			BaseAggregate: evol.NewBaseAggregate(PaymentAggregateType, s),
			PaymentId:     s, //PaymentId as Aggregate root identity
		}
	})
}

var PaymentAggregateType evol.AggregateType = "PaymentAggregate"

type PaymentAggregate struct {
	*evol.BaseAggregate
	PaymentId string //aggregate root identity
	OrderId   string
	Status    string
	Amount    float32
}

func (p *PaymentAggregate) HandleCommand(ctx context.Context, c evol.Command) error {
	switch cmd := c.(type) {
	case *PayOrderCmd:
		buyerBalance := Balance[cmd.BuyerId] //mock, should load form db or query service
		if buyerBalance >= cmd.Amount {
			p.PublishEvent(OrderPayedEventTopic, &OrderPayedEvent{
				PaymentId: cmd.PaymentId,
				OrderId:   cmd.OrderId,
				Amount:    cmd.Amount,
				Time:      time.Now(),
			}, time.Now(), p)
		} else {
			p.PublishEvent(OrderPayFailedEventTopic, &OrderPayFailedEvent{
				PaymentId: cmd.PaymentId,
				OrderId:   cmd.OrderId,
				Reason:    "balance not enough",
				Time:      time.Now(),
			}, time.Now(), p)
		}

	}
	return nil
}

func (p *PaymentAggregate) HandleSourcingEvent(ctx context.Context, e evol.Event) error {
	switch e.Topic() {
	case OrderPayedEventTopic:
		evt := new(OrderPayedEvent)
		err := mapstructure.Decode(e.Data(), evt)
		if err != nil {
			return err
		}

		p.PaymentId = evt.PaymentId
		p.OrderId = evt.OrderId
		p.Amount = evt.Amount
		p.Status = "PAYED"
	case OrderPayFailedEventTopic:
		evt := new(OrderPayFailedEvent)
		err := mapstructure.Decode(e.Data(), evt)
		if err != nil {
			return err
		}

		p.Status = "Failed"
		p.PaymentId = evt.PaymentId
		p.OrderId = evt.OrderId
	}
	return nil
}
