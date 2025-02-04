package api

import (
	"ecommerce/Order/Domain/models"
	"ecommerce/Order/Domain/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/place", h.handlePlaceOrder)
	router.GET("/", h.handleGetOrder)
	router.POST("/{orderId}/cancel", h.handleCancelOrder)
}

func (h *OrderHandler) Name() string {
	return "order"
}

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

	if err := h.orderService.PlaceOrder(request.CustomerID, request.Items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order placed successfully",
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

	if err := h.orderService.CancelOrder(orderId, request.Reason); err != nil {
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
	var request struct {
		OrderId string `json:"order_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	orderId, err := models.NewOrderID(request.OrderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid OrderID: " + err.Error(),
		})
		return
	}

	order, err := h.orderService.GetOrderById(orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Order": order,
	})
}
