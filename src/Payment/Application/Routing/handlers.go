package application

import (
	"ecommerce/Payment/Domain/models"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *PaymentHandler) handleGetPaymentByOrderId(c *gin.Context) {
	orderId := c.Param("orderId")

	payment, err := h.PaymentService.GetPaymentByOrderId(c, orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment": payment,
	})
}

func (h *PaymentHandler) handleGetPayment(c *gin.Context) {
	paymentId := parseToDomainIdAndReturnIfInvalid(c)
	payment, err := h.PaymentService.GetPaymentById(c, paymentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment": payment,
	})
}

func (h *PaymentHandler) handleCompletePayment(c *gin.Context) {
	paymentId := parseToDomainIdAndReturnIfInvalid(c)

	if err := h.PaymentService.ConfirmPayment(c, paymentId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment completed successfully",
	})
}

func (h *PaymentHandler) handleCancelPayment(c *gin.Context) {
	paymentId := parseToDomainIdAndReturnIfInvalid(c)

	if err := h.PaymentService.CancelPayment(c, paymentId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment canceled successfully",
	})
}

func parseToDomainIdAndReturnIfInvalid(c *gin.Context) models.PaymentID {
	paymentId, err := uuid.Parse(c.Param("paymentId"))
	if err != nil {
		slog.Error("Error converting paymentID to uuid", slog.Attr{Key: "error", Value: slog.AnyValue(err)})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	return paymentId
}
