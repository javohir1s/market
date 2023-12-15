package handler

import (
	"errors"
	"net/http"

	"market_system/pkg/security"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		bearerToken := c.GetHeader("Authorization")
		if len(bearerToken) <= 0 {
			c.AbortWithError(http.StatusUnauthorized, errors.New("User not authentication"))
			return
		}

		token, err := security.ExtractToken(cast.ToString(bearerToken))
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		authInfo, err := security.ParseClaims(token, h.cfg.SecretKey)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.Set("user_id", authInfo["user_id"])
		c.Set("client_type", authInfo["client_type"])

		c.Next()
	}
}
