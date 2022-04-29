package adapter

import (
	"context"
	"evol/example/application"
	"evol/example/request"
	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var r request.CreateOrderRequest
	c.Bind(&r)
	application.CreateOrder(context.Background(), &r)

	c.JSON(200, gin.H{
		"message": "success",
	})

}

func CancelOrder(c *gin.Context) {
	var r request.CancelOrderRequest
	c.Bind(&r)
	application.CancelOrder(context.Background(), &r)

	c.JSON(200, gin.H{
		"message": "success",
	})
}

func PayOrder(c *gin.Context) {
	var r request.PayOrderRequest
	c.Bind(&r)
	application.PayOrder(context.Background(), &r)

	c.JSON(200, gin.H{
		"message": "success",
	})
}

func QueryOrders(c *gin.Context) {
	uid := c.Param("uid")
	c.JSON(200, gin.H{"orders": application.QueryOrders(context.Background(), uid)})
}
