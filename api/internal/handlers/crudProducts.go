package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const productServiceURL = "http://products-svc:8081"

func GetProducts(c *gin.Context) {

	resp, err := Chttp.Get(productServiceURL + "/products")

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "product-service unreachable: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read response"})
		return
	}

	c.Data(resp.StatusCode, ContentType, body)

}

func GetProductId(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}

	resp, err := Chttp.Get(productServiceURL + "/products/" + id)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "product-service unreachable: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read response"})
		return
	}

	c.Data(resp.StatusCode, ContentType, body)

}
