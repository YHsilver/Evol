package main

import (
	"context"
	"evol/aggregatestore"
	"evol/application"
	"evol/command"
	"evol/eventbus/local"
	"evol/example/adapter"
	"evol/repo/memory"
	"evol/saga"
	"github.com/gin-gonic/gin"
)

func main() {
	initEvol()
	RunGin()
}

func initEvol() {
	cmdBus := command.NewCommandBus()
	evtBus := local.NewEventBus()
	aggStore := aggregatestore.NewAggregateEventStore(memory.NewAggregateEventRepo())
	sagaStore := saga.NewMemorySagaRepo()
	err := application.Run(context.Background(), cmdBus, evtBus, aggStore, sagaStore)
	if err != nil {
		panic(err)
	}

}

func RunGin() {
	r := gin.Default()
	r.POST("/createOrder", adapter.CreateOrder)
	r.POST("/cancelOrder", adapter.CancelOrder)
	r.POST("/payOrder/:uid", adapter.PayOrder)

	r.GET("/orders", adapter.QueryOrders)

	r.Run()
}
