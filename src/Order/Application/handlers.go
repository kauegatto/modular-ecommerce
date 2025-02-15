package application

import (
	dto "ecommerce/Order/Application/Dto"
	"ecommerce/Order/Domain/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *OrderHandler) handlePlaceOrder(c *gin.Context) {
	var request struct {
		CustomerID string             `json:"customer_id"`
		Items      []models.OrderItem `json:"items"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	order, err := h.orderService.PlaceOrder(c, request.CustomerID, request.Items)
	orderID := order.ID.String()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	orderURL := fmt.Sprintf("http://%s/api/order/%s", c.Request.Host, orderID)

	log.Printf("Order created: ID=%s, CustomerID=%s", orderID, request.CustomerID)

	c.Header("Location", orderURL)
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Order placed successfully",
		"order_id": orderID,
		"link":     orderURL,
	})
}

func (h *OrderHandler) handleCancelOrder(c *gin.Context) {
	var request struct {
		Reason string `json:"reason"`
	}

	orderIDString := c.Param("orderId")

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	orderId, err := models.NewOrderID(orderIDString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid OrderID: " + err.Error(),
		})
		return
	}

	if err := h.orderService.CancelOrder(c, orderId, request.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order canceled successfully",
	})
}

func (h *OrderHandler) handleGetOrder(c *gin.Context) {
	orderIDString := c.Param("orderId")
	h.log.Printf("Getting order with id %s", orderIDString)

	orderId, err := models.NewOrderID(orderIDString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid OrderID: " + err.Error(),
		})
		return
	}

	order, err := h.orderService.GetOrderById(c, orderId)
	orderDto := dto.ToOrderDto(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Order": orderDto,
	})
}
