package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/azuar4e/product/internal/handlers"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	pb "github.com/azuar4e/proto"
)

type Product struct {
	Id    string
	Name  string
	Price int64
}

var products = []Product{
	{Id: "1", Name: "Mechanic Keyboard", Price: 100},
	{Id: "2", Name: "2K Monitor", Price: 300},
	{Id: "3", Name: "Mouse", Price: 50},
}

type server struct {
	pb.UnimplementedProductServiceServer
}

func (s *server) GetProducts(ctx context.Context, in *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	res := make([]*pb.Product, len(products))

	for i := range products {
		res[i] = &pb.Product{
			Id:    products[i].Id,
			Name:  products[i].Name,
			Price: products[i].Price,
		}
	}

	return &pb.GetProductsResponse{Pr: res}, nil
}

func (s *server) GetProductById(ctx context.Context, in *pb.GetProductByIdRequest) (*pb.GetProductByIdResponse, error) {
	for i := range products {
		if products[i].Id == in.Id {
			return &pb.GetProductByIdResponse{
				Pr: &pb.Product{
					Id:    products[i].Id,
					Name:  products[i].Name,
					Price: products[i].Price,
				},
			}, nil
		}
	}
	return nil, errors.New("Invalid ID!")
}

func main() {
	// code of the rest server
	r := gin.Default()
	r.GET("/products", handlers.GetProducts)
	r.GET("/products/:id", handlers.GetProductById)
	fmt.Println("Listening on port 8081")
	go r.Run(":8081")

	// code of the grpc server
	lis, err := net.Listen("tcp", ":8077")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterProductServiceServer(s, &server{})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
