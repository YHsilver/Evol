package domain

// mock DB
var (
	Balance = make(map[string]float32)
	Stock   = make(map[string]int)
	Orders  = make(map[string][]Order, 0)
)

type Order struct {
	OrderId    string
	BuyerId    string
	PaymentId  string
	TotalPrice float32
	ProductIds []string
	Status     string
}

func init() {
	Balance["u100"] = 1000
	Stock["p100"] = 10
}
