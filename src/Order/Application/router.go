package application

import (
	"log"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *OrderService
	log          log.Logger
}

func NewOrderHandler(orderService *OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		log:          *log.Default(),
	}
}

func (h *OrderHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/place", h.handlePlaceOrder)
	router.GET("/:orderId", h.handleGetOrder)
	router.POST("/:orderId/cancel", h.handleCancelOrder)
}

func (h *OrderHandler) Name() string {
	return "order"
}
