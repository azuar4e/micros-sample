package main

import (
	"fmt"

	"github.com/azuar4e/product/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/products", handlers.GetProducts)
	r.GET("/products/:id", handlers.GetProductById)
	r.Run(":8081")
	fmt.Println("Listening on port 8081")
}
