package token

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/montinger-com/montinger-server/app/shared/models/response_model"
	"github.com/montinger-com/montinger-server/app/utils/helpers"
	"github.com/montinger-com/montinger-server/lib/exceptions"
	jwt_utils "github.com/montinger-com/montinger-server/lib/jwt"
)

func Interceptor() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method == "OPTIONS" || c.Request.URL.Path == "/auth/login" || strings.HasPrefix(c.Request.URL.Path, "/metrics") ||
			c.Request.URL.Path == "/monitors/register" || strings.HasPrefix(c.Request.URL.Path, "/monitors/") && strings.HasSuffix(c.Request.URL.Path, "/push") {
			c.Next()
			return
		}

		token := helpers.GetAuthToken(c)

		if token == "" {
			c.JSON(http.StatusUnauthorized, response_model.Result{Message: exceptions.InvalidToken.Error()})
			c.Abort()
			return
		}

		var claims *jwt.Token
		var err error

		if c.Request.URL.Path == "/auth/refresh" {
			claims, err = jwt_utils.ValidateRefreshToken(token)
		} else {
			claims, err = jwt_utils.ValidateAccessToken(token)
		}

		if err != nil {
			c.JSON(http.StatusUnauthorized, response_model.Result{Message: exceptions.InvalidToken.Error(), Errors: []string{err.Error()}})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
