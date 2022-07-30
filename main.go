package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"pilipili/util"
	"syscall"
	"time"

	"pilipili/conf"
	"pilipili/server"
)

func main() {
	// 从配置文件读取配置
	conf.Init()

	// 装载路由
	r := server.NewRouter()
	bind := os.Getenv("API_BIND")

	srv := &http.Server{
		Addr:    bind,
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			util.Log().Panic("listen: %s\n", err)
		}
	}()

	// graceful shutdown implement，copied from gin office docs
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	util.Log().Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		util.Log().Panic("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		util.Log().Println("timeout of 2 seconds.")
	}
	util.Log().Println("Server exiting")
}
