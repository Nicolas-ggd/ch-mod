package handler

import (
	"github.com/Nicolas-ggd/ch-mod/internal/db/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct{}

func (h *Handler) GetUserClaims(c *gin.Context) *models.TokenClaim {
	userClaims, exists := c.Get("user_claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return nil
	}
	return userClaims.(*models.TokenClaim)
}
