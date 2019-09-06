package tasks

import (
	"log"
	"reflect"
	"runtime"
	"time"

	"github.com/robfig/cron"
)

// Cron 定时器单例
var Cron *cron.Cron

// Run 运行
func Run(job func() error) {
	from := time.Now().UnixNano()
	err := job()
	to := time.Now().UnixNano()
	jobName := runtime.FuncForPC(reflect.ValueOf(job).Pointer()).Name()
	if err != nil {
		log.Printf("%s error: %dms\n", jobName, (to-from)/int64(time.Millisecond))
	} else {
		log.Printf("%s success: %dms\n", jobName, (to-from)/int64(time.Millisecond))
	}
}

// CronJob 定时任务
func CronJob() {
	if Cron == nil {
		Cron = cron.New()
	}

	// 注意 spec 参数，这里设置的是 单日零点清空
	Cron.AddFunc("0 0 0 * * *", func() { Run(ResetDailyRank) })
	Cron.Start()

	log.Println("[LOG] Cron job start.....")
}
