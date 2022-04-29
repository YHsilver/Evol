package application

import (
	"context"
	"evol"
	"evol/example/request"

	"evol/example/domain"
)

func CreateOrder(ctx context.Context, req *request.CreateOrderRequest) {
	var totalPrice float32 = 0.0
	for _, id := range req.ProductIds {
		totalPrice += GetPrice(id)
	}

	cmd := &domain.CreateOrderCmd{
		OrderId: req.OrderId,
		BuyerId: req.BuyerId,
		Price:   totalPrice,
		Goods:   req.ProductIds,
	}

	evol.SendCommand(ctx, cmd)
}

func CancelOrder(ctx context.Context, req *request.CancelOrderRequest) {
	//TODO:
}

func PayOrder(ctx context.Context, req *request.PayOrderRequest) {
	//TODO:
}

func QueryOrders(ctx context.Context, uid string) Orders {
	orders := domain.QueryUserOrders(uid)
	return Orders{
		orders: orders,
	}
}

func GetPrice(pid string) float32 {
	return 100
}

type Orders struct {
	orders []domain.Order
}
