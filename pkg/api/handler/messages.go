package handler

import (
	"fmt"
	"github.com/Nicolas-ggd/ch-mod/internal/db/models/request"
	"github.com/Nicolas-ggd/ch-mod/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MessagesHandler struct {
	*Handler
	ChatService services.IChatService
	AuthService services.IAuthService
}

func NewMessageHandler(chatService services.IChatService, authService services.IAuthService) *MessagesHandler {
	return &MessagesHandler{
		ChatService: chatService,
		AuthService: authService,
	}
}

func (mh *MessagesHandler) SendMessage(c *gin.Context) {
	var message request.ChatRequest

	err := c.ShouldBind(&message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (mh *MessagesHandler) SendMessageList(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	usr := mh.GetUserClaims(c)

	fmt.Println(usr.UserId)

	model, err := mh.ChatService.FindByUser(uint(userId), usr.UserId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": model})
}

func (mh *MessagesHandler) UserConversations(c *gin.Context) {
	usr := mh.GetUserClaims(c)

	model, err := mh.ChatService.UserConversations(usr.UserId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": model})
}
