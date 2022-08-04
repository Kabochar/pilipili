package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"pilipili/conf"
	"pilipili/server"
	"pilipili/util"
)

func main() {
	// 从配置文件读取配置
	conf.Init()

	// 装载路由
	r := server.NewRouter()

	// 获取bind地址和当前服务名称
	var addr string = os.Getenv("API_BIND")
	var svcName string = os.Getenv("SERVICE_NAME")

	// 创建/导入监听器
	ln, err := util.CreateOrImportListener(svcName, addr)
	if err != nil {
		fmt.Printf("Unable to create or import a listener: %v.\n", err)
		return
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	go func() {
		// service connections
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			util.Log().Panic("listen: %s\n", err)
		}
	}()

	// 接收退出信号
	err = util.WaitForSignals(svcName, addr, ln, srv)
	if err != nil {
		util.Log().Error("Exiting: %v\n", err)
		return
	}
}
