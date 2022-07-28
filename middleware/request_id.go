package middleware

import (
	"pilipili/util"

	"github.com/gin-gonic/gin"
)

// 请求id
// 参考链接：https://mp.weixin.qq.com/s/M2jNnLkYaearwyRERnt0tA
func RequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.Request.Header.Get("X-Trace-Id")
		if traceID == "" {
			traceID = util.RandStringRunes(15)
			ctx.Request.Header.Set("X-Trace-Id", traceID)
		}
		spanID := ctx.Request.Header.Get("X-Span-Id")
		if spanID == "" {
			spanID = getSpanID(ctx.RemoteIP())
			ctx.Request.Header.Set("X-Span-Id", spanID)
		}

		// 写入响应
		ctx.Header("X-Trace-Id", traceID)
		ctx.Header("X-Span-Id", spanID)
		ctx.Next()
	}
}

// 根据客户端ip生成唯一
func getSpanID(ip string) string {
	result := util.MD5(ip)
	return result[15:]
}
