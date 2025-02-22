package application

import (
	service "ecommerce/Payment/Application/Service"
	"log"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	PaymentService *service.PaymentService
	log            log.Logger
}

func NewPaymentHandler(PaymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		PaymentService: PaymentService,
		log:            *log.Default(),
	}
}

func (h *PaymentHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/order/:orderId", h.handleGetPaymentByOrderId)
	router.GET("/:paymentId", h.handleGetPayment)

	router.POST("/:paymentId/capture", h.handleCapturePayment)
	router.POST("/:paymentId/complete", h.handleCompletePayment)
}

func (h *PaymentHandler) Name() string {
	return "payment"
}
