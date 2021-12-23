package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	conf "price_api/price_server/config"
	"price_api/price_server/internal/constant"
	"price_api/price_server/internal/pkg/jwt"
	"price_api/price_server/internal/vo"
	"strings"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		response := vo.RESPONSE{Code: 0, Message: "OK"}

		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			response.Code = constant.CHECK_USER_ERROR
			response.Message = constant.MSG_CHECK_USER_ERROR
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Code = constant.CHECK_USER_ERROR
			response.Message = constant.MSG_CHECK_USER_ERROR
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		mc, err := jwt.ParseToken(parts[1], []byte(conf.GCfg.Password))
		if err != nil {
			response.Code = constant.CHECK_USER_ERROR
			response.Message = constant.MSG_CHECK_USER_ERROR
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		c.Set("username", mc.Username)
		c.Next()
	}
}
