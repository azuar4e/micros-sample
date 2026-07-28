package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/azuar4e/api/internal/auth"
	"github.com/azuar4e/api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	if err := auth.InitAuth(ctx); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	v1 := r.Group("/api/v1")

	// authentications against idp
	v1.GET("/authentication", auth.Login)
	v1.GET("/callback", auth.Callback)

	v1.Use(auth.AuthMiddleware(), auth.AuthzMiddleware("customer"))

	// protected uris
	v1.GET("/me", func(c *gin.Context) {
		userID, _ := c.Get("userID")
		email, _ := c.Get("userEmail")
		roles, _ := c.Get("roles")

		c.JSON(http.StatusOK, gin.H{
			"message": "I'm logged in!",
			"userID":  userID,
			"email":   email,
			"roles":   roles,
		})
	})

	v1.GET("/products", handlers.GetProducts)
	v1.GET("/products/:id", handlers.GetProductId)

	v1.POST("/orders", handlers.PostOrder)
	v1.GET("/orders/:id", handlers.GetOrderId)

	r.Run(":9090")
	fmt.Println("Listening on port 9090")
}
