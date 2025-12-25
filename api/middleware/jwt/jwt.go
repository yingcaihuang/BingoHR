package jwt

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"hr-api/pkg/namedroute"
	"hr-api/pkg/util"
	"hr-api/service/user_service"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var msg string
		var data interface{}

		token := c.GetHeader("x-token")
		if token == "" {
			code = 1
			msg = "缺少x-token鉴权头"
		} else {
			claim, err := util.ParseToken(token)
			if err != nil {
				code = 1
				// v5 版本错误处理
				if errors.Is(err, jwt.ErrTokenExpired) {
					msg = "Token已超时"
				} else {
					msg = "Token鉴权失败"
				}
			} else {
				c.Set("uid", claim.Id)
				c.Set("username", claim.Username)

				if claim.Id > 1 {
					routeName := namedroute.GetRouteName(c)
					service := user_service.User{Id: claim.Id}
					perms := service.GetUserPerms()
					if len(routeName) > 0 && len(perms) > 0 && !util.Contains(perms, routeName) {
						c.JSON(http.StatusForbidden, gin.H{
							"code": 1,
							"msg":  "Permission deined",
						})
						c.Abort()
						return
					}
				}
			}
		}

		if code > 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  msg,
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
