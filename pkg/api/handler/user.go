package handler

import (
	"github.com/Nicolas-ggd/ch-mod/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	*Handler
	UserService services.IUserService
	AuthService services.IAuthService
}

func NewUserHandler(userService services.IUserService, authService services.IAuthService) *UserHandler {
	return &UserHandler{
		UserService: userService,
		AuthService: authService,
	}
}

func (h *UserHandler) FindByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	user, err := h.UserService.FindByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Me(c *gin.Context) {
	userToken := h.GetUserClaims(c)

	user, err := h.UserService.FindByID(userToken.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
