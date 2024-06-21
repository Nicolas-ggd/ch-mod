package middleware

import (
	"github.com/Nicolas-ggd/ch-mod/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(authService services.IAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {

		userObj, err := authService.CheckJWT(c.GetHeader("Authorization"))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}

		c.Set("user_claims", userObj)

		c.Next()
	}

}
