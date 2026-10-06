package handlers

import (
	"context"
	"net/http"
	"sync"
	"time"

	pb "github.com/azuar4e/proto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Order struct {
	Id        string `json:"id"`
	ProductId string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

var orders sync.Map

type OrderHandler struct {
	Products pb.ProductServiceClient
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	// Timeout for context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var reqOrder struct {
		ProductId string `json:"product_id" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&reqOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.Products.GetProductById(ctx, &pb.GetProductByIdRequest{
		Id: reqOrder.ProductId,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	orderId := uuid.New().String()
	order := Order{Id: orderId, ProductId: reqOrder.ProductId, Quantity: reqOrder.Quantity}
	orders.Store(orderId, order)
	c.JSON(http.StatusOK, gin.H{"data": order})
}

func (h *OrderHandler) GetOrderById(c *gin.Context) {
	id := c.Param("id")

	orderData, ok := orders.Load(id)

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	order, ok := orderData.(Order)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "corrupted order data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}
