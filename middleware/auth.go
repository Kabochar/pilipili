package middleware

import (
	"pilipili/model"
	"pilipili/serializer"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// CurrentUser 获取登录用户
func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		uid := session.Get("user_id")
		if uid == nil {
			c.Next()
			return
		}

		user, err := model.GetUser(uid)
		if err != nil {
			c.Next()
			return
		}
		c.Set("user", &user)
		c.Next()
	}
}

// AuthRequired 需要登录
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("user")
		if user == nil {
			c.JSON(200, serializer.Response{
				Status: 401,
				Msg:    "需要登录",
			})
			c.Abort()
		}
		if _, ok := user.(*model.User); ok {
			c.Next()
		}
	}
}
