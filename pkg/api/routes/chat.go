package routes

import (
	"github.com/Nicolas-ggd/ch-mod/pkg/api/handler"
	"github.com/Nicolas-ggd/ch-mod/pkg/api/middleware"
	"github.com/gin-gonic/gin"
)

func ChatRoutes(e *gin.RouterGroup, chatHandler *handler.MessagesHandler) {
	e.Use(middleware.AuthMiddleware(chatHandler.AuthService))
	e.GET("/chat-list/:id", chatHandler.SendMessageList)
	e.GET("/user-conversation", chatHandler.UserConversations)
}
