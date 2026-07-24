package main

import (
	"fmt"

	"github.com/azuar4e/order/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// order-service
	r := gin.Default()
	r.POST("/orders", handlers.CreateOrder)
	r.GET("/orders/:id", handlers.GetOrderById)
	r.Run(":8082")
	fmt.Println("Listening on port 8082")
}
