package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	router *gin.Engine
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", h.handleX)
	router.POST("", h.handleY)
}

func (h *AuthHandler) Name() string {
	return "auth"
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) handleX(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *AuthHandler) handleY(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
