package api

import (
	"github.com/Nicolas-ggd/ch-mod/pkg/api/handler"
	"github.com/Nicolas-ggd/ch-mod/pkg/api/routes"
	"github.com/Nicolas-ggd/ch-mod/pkg/api/ws"
	"github.com/Nicolas-ggd/ch-mod/pkg/repository"
	"github.com/Nicolas-ggd/ch-mod/pkg/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ServeAPI(db *gorm.DB) *gin.Engine {
	repositories := repository.NewRepository(db)
	service := services.NewService(repositories)

	authHandler := handler.NewAuthHandler(service.AuthService)

	r := gin.Default()

	wss := ws.NewWebsocket(service.ChatService)

	go wss.Run()

	r.GET("/ws", wss.ServeWs)
	v1 := r.Group("v1")
	{
		routes.AuthRoutes(v1.Group("auth"), authHandler)
	}

	return r
}
