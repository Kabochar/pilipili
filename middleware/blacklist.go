package middleware

import (
	"net/http"

	"pilipili/model"
	"pilipili/serializer"
	"pilipili/util"

	"github.com/gin-gonic/gin"
)

func BlacklistMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("user")
		if user == nil {
			c.Next()
			return
		}
		userInfo, ok := user.(*model.User)
		if !ok {
			c.Next()
			return
		}
		if util.CheckInBlackList(userInfo.ID) {
			c.JSON(200, serializer.Response{
				Status: http.StatusForbidden,
				Msg:    "禁止访问",
			})
			c.Abort()
		}
		c.Next()
	}
}
