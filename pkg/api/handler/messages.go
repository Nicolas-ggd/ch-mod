package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type MessagesHandler struct {
	*Handler
}

func (mh *MessagesHandler) SendMessage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
