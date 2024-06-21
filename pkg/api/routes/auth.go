package routes

import (
	"github.com/Nicolas-ggd/ch-mod/pkg/api/handler"
	"github.com/Nicolas-ggd/ch-mod/pkg/api/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(e *gin.RouterGroup, authHandler *handler.AuthHandler) {
	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)
	e.Use(middleware.AuthMiddleware(authHandler.AuthService))
	e.POST("/logout", authHandler.Logout)
}
