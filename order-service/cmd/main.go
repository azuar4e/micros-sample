package main

import (
	"fmt"
	"log"

	"github.com/azuar4e/order/internal/handlers"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/azuar4e/proto"
)

func main() {
	// order grpc client
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("products-svc:8077", opts...)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := &handlers.OrderHandler{Products: pb.NewProductServiceClient(conn)}

	// order-service
	r := gin.Default()
	r.POST("/orders", client.CreateOrder)
	r.GET("/orders/:id", client.GetOrderById)

	fmt.Println("Listening on port 8082")
	log.Fatal(r.Run(":8082"))
}
