package model

import (
	"log"
	"os"
	"time"

	"pilipili/util"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 数据库链接单例
var DB *gorm.DB

// Database 在中间件中初始化mysql链接
func Database(connString string) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	db, err := gorm.Open(sqlite.Open(os.Getenv("SQLITE_DB")), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		util.Log().Panic("connect mysql client ERROR %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		util.Log().Panic("get mysql client ERROR %v", err)
	}
	// 设置连接池
	// 空闲
	sqlDB.SetMaxIdleConns(50)
	// 打开
	sqlDB.SetMaxOpenConns(100)
	// 超时
	sqlDB.SetConnMaxLifetime(time.Second * 30)

	DB = db

	// 执行数据迁移
	migration()
}
