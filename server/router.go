package server

import (
	"os"

	"github.com/gin-gonic/gin"
	"pilipili/api"
	"pilipili/middleware"
	"pilipili/util"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	gin.DisableConsoleColor()

	r := gin.Default()

	// 中间件, 顺序不能改
	r.Use(middleware.Session(os.Getenv("SESSION_SECRET"))) // 这里干嘛，加密 session，创建一个唯一性的 session
	r.Use(middleware.Cors())                               // 跨域问题
	r.Use(middleware.CurrentUser())                        // 获取当前用户

	// 配置可信任的代理，配置为 nil，默认都允许通过
	// 参考资料：https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	err := r.SetTrustedProxies(nil)
	if err != nil {
		util.Log().Error("set trusted proxies ERROR [%v]\n", err)
	}

	// 路由
	v1 := r.Group("/api/v1")
	{
		// 服务存活判断
		v1.POST("ping", api.Ping)

		// 用户登录
		v1.POST("user/register", api.UserRegister)

		// 用户登录
		v1.POST("user/login", api.UserLogin)

		// 需要登录保护的
		auth := v1.Group("/")
		auth.Use(middleware.AuthRequired()) // 身份认证
		{
			// User Routing
			auth.GET("user/me", api.UserMe)
			auth.DELETE("user/logout", api.UserLogout)
		}

		// 视频相关服务
		v1.POST("videos", api.CreateVideo)
		v1.GET("video/:id", api.ShowVideo)
		v1.GET("videos", api.ListVideo)
		v1.PUT("video/:id", api.UpdateVideo)
		v1.DELETE("video/:id", api.DeleteVideo)

		// 排行榜
		v1.GET("rank/daily", api.DailyRank)

		// 文件上传
		v1.POST("upload/token", api.UploadToken)
	}

	return r
}
