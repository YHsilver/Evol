package request

type CreateOrderRequest struct {
	OrderId    string   `form:"orderId"`
	BuyerId    string   `form:"userId"`
	ProductIds []string `form:"products"`
}

type CancelOrderRequest struct {
	OrderId string `form:"orderId"`
	BuyerId string `form:"userId"`
}

type PayOrderRequest struct {
	OrderId string `form:"orderId"`
	BuyerId string `form:"userId"`
	PayType string `form:"payType"`
}
