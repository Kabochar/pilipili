package main

import (
	"net/http"
	"os"
	"time"

	"pilipili/conf"
	"pilipili/server"
	"pilipili/util"
)

func main() {
	// 从配置文件读取配置
	conf.Init()

	// 装载路由
	router := server.NewRouter()

	// 获取bind地址和当前服务名称
	var addr string = os.Getenv("API_BIND")
	var svcName string = os.Getenv("SERVICE_NAME")

	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 2 * time.Second,
	}
	// 开始监听服务
	util.GracefulShutdown(svcName, addr, srv)
}
