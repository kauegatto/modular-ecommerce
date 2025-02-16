package application

import (
	"github.com/gin-gonic/gin"
)

func (h *PaymentHandler) handleGetPayment(c *gin.Context) {
	c.JSON(200, gin.H{
		"response": "ok",
	})
}

// those two would probably be webhooks and contain keys to process operations

func (h *PaymentHandler) handleCompletePayment(c *gin.Context) {
	c.JSON(200, gin.H{
		"response": "ok",
	})
}

func (h *PaymentHandler) handleCancelPayment(c *gin.Context) {
	c.JSON(200, gin.H{
		"response": "ok",
	})
}
