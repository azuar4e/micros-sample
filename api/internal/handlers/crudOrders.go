package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const ordersServiceURL = "http://orders-svc:8082"
const ContentType = "application/json"

func PostOrder(c *gin.Context) {

	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	// Inyectas el user_id del token para que order-service sepa de quién es el pedido
	userID, _ := c.Get("userID")
	reqBody["user_id"] = userID

	jsonBody, _ := json.Marshal(reqBody)

	resp, err := Chttp.Post(ordersServiceURL+"/orders", ContentType, bytes.NewReader(jsonBody))

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "orders-service unreachable: " + err.Error()})
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

func GetOrderId(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}

	resp, err := Chttp.Get(ordersServiceURL + "/orders/" + id)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "orders-service unreachable: " + err.Error()})
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
