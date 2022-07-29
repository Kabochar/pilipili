package middleware

import (
	"pilipili/util"

	"github.com/gin-gonic/gin"
)

const (
	TRACK_ID_HEADER = "X-Trace-Id" // trace 串联出一个业务逻辑的所有业务日志
	SPAN_ID_HEADER  = "X-Span-Id"  // span 串联出在单个服务里的业务日志
)

// 请求id
// 参考链接：https://mp.weixin.qq.com/s/M2jNnLkYaearwyRERnt0tA
func RequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.Request.Header.Get(TRACK_ID_HEADER)
		spanID := getSpanID(ctx.ClientIP())
		if traceID == "" {
			traceID = spanID
			// 写入请求header
			ctx.Request.Header.Set(TRACK_ID_HEADER, traceID)
			ctx.Request.Header.Set(SPAN_ID_HEADER, spanID)
		}

		// 写入响应header
		ctx.Header(TRACK_ID_HEADER, traceID)
		ctx.Header(SPAN_ID_HEADER, spanID)
		ctx.Next()
	}
}

// 根据客户端ip生成唯一
func getSpanID(ip string) string {
	result := util.MD5(ip)
	return result[15:]
}
