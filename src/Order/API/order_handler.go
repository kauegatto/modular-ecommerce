package api

import (
	domain "ecommerce/Order/Domain"
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
}

func (h *OrderHandler) Name() string {
	return "order"
}

func (h *OrderHandler) handlePlaceOrder(c *gin.Context) {
	var request struct {
		CustomerID string  `json:"customer_id"`
		Amount     float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	if err := h.orderService.PlaceOrder(request.CustomerID, domain.Money(request.Amount)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order placed successfully",
	})
}
