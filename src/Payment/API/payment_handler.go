package api

import (
	"ecommerce/Payment/Domain/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	PaymentService *services.PaymentService
}

func NewPaymentHandler(PaymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		PaymentService: PaymentService,
	}
}

func (h *PaymentHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/place", h.handlePlacePayment)
}

func (h *PaymentHandler) Name() string {
	return "Payment"
}

func (h *PaymentHandler) handlePlacePayment(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Payment placed successfully",
	})
}
