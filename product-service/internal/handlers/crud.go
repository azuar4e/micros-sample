package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

var products = []Product{
	{Id: "1", Name: "Mechanic Keyboard", Price: 100},
	{Id: "2", Name: "2K Monitor", Price: 300},
	{Id: "3", Name: "Mouse", Price: 50},
}

func GetProducts(c *gin.Context) {
	c.JSON(http.StatusOK, products)
}

func GetProductById(c *gin.Context) {
	id := c.Param("id")

	for _, product := range products {
		if product.Id == id {
			c.JSON(http.StatusOK, product)
			return
		}
	}
}
