package util

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// LevelError 错误
	LevelError = iota
	// LevelWarning 警告
	LevelWarning
	// LevelInformational 提示
	LevelInformational
	// LevelDebug 出错
	LevelDebug
)

var logger *Logger

// Logger 日志
type Logger struct {
	level   int
	traceID string // trace 串联出一个业务逻辑的所有业务日志
	spanID  string // span 串联出在单个服务里的业务日志
}

// Println 打印
func (ll *Logger) Println(msg string) {
	fmt.Printf("[%s %d %s %s] %s\n",
		time.Now().Format("2006-01-02 15:04:05.000"),
		os.Getpid(),
		ll.traceID,
		ll.spanID,
		msg,
	)
}

// Panic 极端错误
func (ll *Logger) Panic(format string, v ...interface{}) {
	if LevelError > ll.level {
		return
	}
	msg := fmt.Sprintf("[Panic] "+format, v...)
	ll.Println(msg)
	os.Exit(0)
}

// Error 错误
func (ll *Logger) Error(format string, v ...interface{}) {
	if LevelError > ll.level {
		return
	}
	msg := fmt.Sprintf("[E] "+format, v...)
	ll.Println(msg)
}

// Warning 警告
func (ll *Logger) Warning(format string, v ...interface{}) {
	if LevelWarning > ll.level {
		return
	}
	msg := fmt.Sprintf("[W] "+format, v...)
	ll.Println(msg)
}

// Info 信息
func (ll *Logger) Info(format string, v ...interface{}) {
	if LevelInformational > ll.level {
		return
	}
	msg := fmt.Sprintf("[I] "+format, v...)
	ll.Println(msg)
}

// Debug 校验
func (ll *Logger) Debug(format string, v ...interface{}) {
	if LevelDebug > ll.level {
		return
	}
	msg := fmt.Sprintf("[D] "+format, v...)
	ll.Println(msg)
}

// BuildLogger 构建logger
func BuildLogger(level string) {
	intLevel := LevelError
	switch level {
	case "error":
		intLevel = LevelError
	case "warning":
		intLevel = LevelWarning
	case "info":
		intLevel = LevelInformational
	case "debug":
		intLevel = LevelDebug
	}
	l := Logger{
		level: intLevel,
	}
	logger = &l
}

// GetTraceId 由中间件传入具体的trace id
func (ll *Logger) GetTraceId(ctx *gin.Context) *Logger {
	ll.traceID = ctx.GetHeader("X-Trace-Id")
	ll.spanID = ctx.GetHeader("X-Span-Id")
	return ll
}

// Log 返回日志对象
func Log() *Logger {
	if logger == nil {
		l := Logger{
			level: LevelDebug,
		}
		logger = &l
	}
	return logger
}
