package middleware

import (
	"pilipili/util"

	"github.com/gin-gonic/gin"
)

// 请求id
func RequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqId := ctx.Request.Header.Get("X-Request-Id")
		if reqId == "" {
			reqId = util.RandStringRunes(15)
			ctx.Request.Header.Set("X-Request-Id", reqId)
		}
		// 写入到service的log对象上
		util.Log().SetLogField(reqId)

		// 写入响应
		ctx.Header("X-Request-Id", reqId)
		ctx.Next()
	}
}
