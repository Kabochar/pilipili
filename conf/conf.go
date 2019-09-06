package conf

import (
	"os"

	"pilipili/cache"
	"pilipili/model"
	"pilipili/tasks"
	"pilipili/util"

	"github.com/joho/godotenv"
)

// Init 初始化配置项
func Init() {
	// 从本地读取环境变量
	godotenv.Load()

	// 设置日志级别
	util.BuildLogger(os.Getenv("LOG_LEVEL"))

	// 连接数据库
	model.Database(os.Getenv("MYSQL_DSN"))
	cache.Redis()

	// 启动定时任务
	tasks.CronJob()
}
