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

	token, userId, err := ah.AuthService.Login(&login)
	if err != nil {
		if strings.Contains(err.Error(), "user email or password is incorrect") {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user_id": userId})
}

func (ah *AuthHandler) Logout(c *gin.Context) {
	userToken := ah.GetUserClaims(c)

	err := ah.AuthService.Logout(userToken.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (ah *AuthHandler) SetPassword(c *gin.Context) {
	var newPassword request.SetPasswordRequest

	err := c.ShouldBind(&newPassword)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	smsCode := c.Query("code")
	if smsCode == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request data, missing code query parameter"})
		return
	}

	err = ah.AuthService.SetPassword(newPassword, smsCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "You set new password"})
}

func (ah *AuthHandler) ResetPassword(c *gin.Context) {
	var email request.ResetPasswordRequest

	err := c.ShouldBind(&email)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	err = ah.AuthService.VerifyCredentials(email.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Check your email"})
}
