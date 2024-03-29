package util

import (
	"time"

	"github.com/chenyahui/gin-cache/persist"
)

const (
	// 默认的缓存中间件过期时间
	_DEFUALT_CACHE_RESPONSE_EXPIRE_TIME = time.Minute
)

var (
	// 响应缓存中间件对象
	RespMemCache *persist.MemoryStore = nil
)

// 构建缓存中间件
func BuildResponseMemoryCache() {
	RespMemCache = persist.NewMemoryStore(_DEFUALT_CACHE_RESPONSE_EXPIRE_TIME)
}
