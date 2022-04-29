package domain

import (
	"context"
	"evol"
	"evol/example/infra"
	"evol/saga"
	"github.com/mitchellh/mapstructure"
	"github.com/thoas/go-funk"
	"strconv"
)

func init() {
	opt := saga.SagaManagerOpt{
		SagaType:    "OrderSaga",
		StartEvents: []evol.Topic{OrderCreatedEventTopic},
		OnEvents:    []evol.Topic{ProductReservedEventTopic, OrderPayedEventTopic},
		EndEvents:   []evol.Topic{ProductReserveFailedEventTopic, OrderPayFailedEventTopic},
		Resolver:    Resolver,
		SagaFactory: NewOrderSaga,
	}
	manager := saga.NewSagaManager(&opt)

	err := saga.RegisterSaga(manager)
	if err != nil {
		panic("Register Saga Manager Failed, err:=" + err.Error())
	}

}

func NewOrderSaga(sagaType string, sagaIdentity string) evol.SagaHandler {
	return &OrderSaga{BaseSaga: saga.NewBaseSaga(sagaIdentity, sagaType)}
}

func Resolver(e evol.Event) string {
	return "OrderId"
}

// OrderSaga handle a transaction across aggregates
// CreateOrderCmd -> OrderAggregate -> OrderCreatedEvent ->
// OrderSaga -> MakeReservationCmd -> StockAggregate -> ProductReservedEvent ->
// OrderSaga -> PayOrderCmd -> PaymentAggregate -> OrderPayedEvent -> OrderSaga.endSaga
type OrderSaga struct {
	*saga.BaseSaga
	OrderId               string
	BuyerId               string
	TotalPrice            float32
	AllProducts           []string
	ReservedProducts      []string
	ReserveFailedProducts []string

	toReserveNum int
	needRollBack bool
}

func (o *OrderSaga) HandleSagaEvent(ctx context.Context, event evol.Event, bus evol.CommandHandler) error {
	switch event.Topic() {
	case OrderCreatedEventTopic: //start saga
		evt := new(OrderCreatedEvent)
		err := mapstructure.Decode(event.Data(), evt)
		if err != nil {
			return err
		}
		o.OrderId = evt.OrderId
		o.AllProducts = evt.ProductIds
		o.BuyerId = evt.BuyerId
		o.TotalPrice = evt.TotalPrice

		for _, productId := range evt.ProductIds {
			cmd := &MakeReservationCmd{
				OrderId:   evt.OrderId,
				ProductId: productId,
				Count:     1,
			}
			//TODO: Rpc send command
			_ = bus.HandleCommand(ctx, cmd)
		}

	case ProductReservedEventTopic:
		evt := new(ProductReservedEvent)
		err := mapstructure.Decode(event.Data(), evt)
		if err != nil {
			return err
		}
		if !funk.Contains(o.ReservedProducts, evt.ProductId) {
			o.ReservedProducts = append(o.ReservedProducts, evt.ProductId)
		}
		o.toReserveNum--
		if o.toReserveNum == 0 {
			return o.FinishReservation(ctx, bus)
		}

	case ProductReserveFailedEventTopic:
		o.needRollBack = true
		evt := new(ProductReserveFailedEvent)
		err := mapstructure.Decode(event.Data(), evt)
		if err != nil {
			return err
		}
		if !funk.Contains(o.ReserveFailedProducts, evt.ProductId) {
			o.ReserveFailedProducts = append(o.ReserveFailedProducts, evt.ProductId)
		}
		o.toReserveNum--
		if o.toReserveNum == 0 {
			return o.FinishReservation(ctx, bus)
		}

	case OrderPayedEventTopic:
		evt := new(OrderPayedEvent)
		err := mapstructure.Decode(event.Data(), evt)
		if err != nil {
			return err
		}
		cmd := &OrderConfirmedCmd{
			OrderId: o.OrderId,
		}
		return bus.HandleCommand(ctx, cmd)

	case OrderPayFailedEventTopic:
		o.needRollBack = true
		evt := new(OrderPayFailedEvent)
		err := mapstructure.Decode(event.Data(), evt)
		if err != nil {
			return err
		}
		o.RollBackOrder(ctx, bus)

		//cancel order
		cmd := &CancelOrderCmd{
			OrderId: o.OrderId,
			Reason:  "Pay Order Failed",
		}
		bus.HandleCommand(ctx, cmd)
	}

	return nil
}

func (o *OrderSaga) FinishReservation(ctx context.Context, bus evol.CommandHandler) error {
	if o.needRollBack {
		o.RollBackOrder(ctx, bus)
		//cancel order
		cmd := &CancelOrderCmd{
			OrderId: o.OrderId,
			Reason:  "Reserve Product Failed",
		}
		bus.HandleCommand(ctx, cmd)
	} else {
		// payOrder
		pid, err := infra.NewUUID()
		if err != nil {
			return err
		}
		cmd := &PayOrderCmd{
			PaymentId: strconv.FormatInt(pid, 10),
			OrderId:   o.OrderId,
			BuyerId:   o.BuyerId,
			Amount:    o.TotalPrice,
		}

		return bus.HandleCommand(ctx, cmd)
	}
	return nil
}

func (o *OrderSaga) RollBackOrder(ctx context.Context, bus evol.CommandHandler) {
	//rollback reserved products
	for _, product := range o.ReservedProducts {
		cmd := &RollBackReservationCmd{
			OrderId:   o.OrderId,
			ProductId: product,
			Count:     1,
		}
		bus.HandleCommand(ctx, cmd)
	}
}
