package handler

import (
	"github.com/Nicolas-ggd/ch-mod/internal/db/models/request"
	"github.com/Nicolas-ggd/ch-mod/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AuthHandler struct {
	*Handler
	AuthService services.IAuthService
}

func NewAuthHandler(authService services.IAuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

func (ah *AuthHandler) Register(c *gin.Context) {
	var register request.UserRegisterRequest

	err := c.ShouldBind(&register)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := ah.AuthService.Register(&register)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (ah *AuthHandler) Login(c *gin.Context) {
	var login request.LoginRequest

	err := c.ShouldBind(&login)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	token, err := ah.AuthService.Login(&login)
	if err != nil {
		if strings.Contains(err.Error(), "user email or password is incorrect") {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (ah *AuthHandler) Logout(c *gin.Context) {

}
