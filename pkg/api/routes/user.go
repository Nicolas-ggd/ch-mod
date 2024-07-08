package routes

import (
	"github.com/Nicolas-ggd/ch-mod/pkg/api/handler"
	"github.com/Nicolas-ggd/ch-mod/pkg/api/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(e *gin.RouterGroup, userHandler *handler.UserHandler) {
	e.Use(middleware.AuthMiddleware(userHandler.AuthService))
	e.GET("/search", userHandler.FindByEmail)
	e.GET("/me", userHandler.Me)
}
