package handlers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Order struct {
	Id        string `json:"id"`
	ProductId string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

var orders sync.Map

func CreateOrder(c *gin.Context) {
	var reqOrder struct {
		ProductId string `json:"product_id" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&reqOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := http.Get("http://products-svc:8081/products/" + reqOrder.ProductId)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	orderId := uuid.New().String()
	order := Order{Id: orderId, ProductId: reqOrder.ProductId, Quantity: reqOrder.Quantity}
	orders.Store(orderId, order)
	c.JSON(http.StatusOK, gin.H{"data": order})
}

func GetOrderById(c *gin.Context) {

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
