package main

import (
	"fmt"

	"github.com/azuar4e/api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	v1 := r.Group("/api/v1")

	v1.GET("/products", handlers.GetProducts)
	v1.GET("/products/:id", handlers.GetProductId)

	v1.POST("/orders", handlers.PostOrder)
	v1.GET("/orders/:id", handlers.GetOrderId)

	r.Run(":9090")
	fmt.Println("Listening on port 9090")
}
