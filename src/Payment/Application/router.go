package application

import (
	"log"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	PaymentService *PaymentService
	log            log.Logger
}

func NewPaymentHandler(PaymentService *PaymentService) *PaymentHandler {
	return &PaymentHandler{
		PaymentService: PaymentService,
		log:            *log.Default(),
	}
}

func (h *PaymentHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/:paymentId", h.handleGetPayment)
	router.POST("/:paymentId/complete", h.handleCompletePayment)
	router.POST("/:paymentId/cancel", h.handleCancelPayment)
}

func (h *PaymentHandler) Name() string {
	return "Payment"
}
