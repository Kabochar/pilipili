package conf

import (
	"io/ioutil"
	"os"

	"pilipili/cache"
	"pilipili/model"
	"pilipili/tasks"
	"pilipili/util"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Init 初始化配置项
func Init() {
	// 从本地读取环境变量
	godotenv.Load()

	// 读取翻译文件
	if err := LoadLocales("conf/locales/zh-cn.yaml"); err != nil {
		panic(err)
	}

	// 设置日志级别
	util.BuildLogger(os.Getenv("LOG_LEVEL"))

	// 配置gin 日志
	gin.SetMode(os.Getenv("GIN_MODE"))
	if os.Getenv("TEST_MODE") == "benchmark" { // 压测模式下，暂时屏蔽日志以测试出理论最优的QPS
		gin.DefaultWriter = ioutil.Discard
	}

	model.Database(os.Getenv("MYSQL_DSN")) // 连接数据库
	cache.Redis()                          // 连接Redis
	tasks.CronJob()                        // 启动定时任务
	util.WatchBlackList()                  // 加载黑名单
	util.BuildResponseMemoryCache()        // 初始化响应中间件对象
}
