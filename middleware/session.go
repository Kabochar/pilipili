package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// Session 初始化session
func Session(secret string) gin.HandlerFunc {
	store := cookie.NewStore([]byte(secret))
	// 设置安全性：如果使用 SSL，放行通过
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 86400, Path: "/"})

	return sessions.Sessions("gin-session", store)
}
