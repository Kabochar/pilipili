package middleware

import (
	"time"

	"github.com/chenyahui/gin-cache/persist"
)

const (
	// 默认的缓存中间件过期时间
	_DEFUALT_CACHE_RESPONSE_EXPIRE_TIME = time.Minute
)

var (
	// 缓存中间件对象
	MemoryCache *persist.MemoryStore = nil
)

// 构建缓存中间件
func BuildMiddleMemoryCache() {
	MemoryCache = persist.NewMemoryStore(_DEFUALT_CACHE_RESPONSE_EXPIRE_TIME)
}
