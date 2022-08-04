package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

/**
Warning(警告)
This page's code copied from this article: https://goteleport.com/blog/golang-ssh-bastion-graceful-restarts/.

For insert this project, may do some transform.
——————————
本代码拷贝自这这篇文章：https://goteleport.com/blog/golang-ssh-bastion-graceful-restarts/.

在符合项目的要求下做了些许改造
*/

type listener struct {
	Addr     string `json:"addr"`
	FD       int    `json:"fd"`
	Filename string `json:"filename"`
}

func getSvcListenKey(svcName, addr string) string {
	return svcName + "_" + addr
}

func importListener(svcName, addr string) (net.Listener, error) {
	// Extract the encoded listener metadata from the environment.
	listenerEnv := os.Getenv(getSvcListenKey(svcName, addr))
	if listenerEnv == "" {
		return nil, fmt.Errorf("unable to find LISTENER environment variable")
	}

	// Unmarshal the listener metadata.
	var l listener
	err := json.Unmarshal([]byte(listenerEnv), &l)
	if err != nil {
		return nil, err
	}
	if l.Addr != addr {
		return nil, fmt.Errorf("unable to find listener for %v", addr)
	}

	// The file has already been passed to this process, extract the file
	// descriptor and name from the metadata to rebuild/find the *os.File for
	// the listener.
	listenerFile := os.NewFile(uintptr(l.FD), l.Filename)
	if listenerFile == nil {
		return nil, fmt.Errorf("unable to create listener file: %v", err)
	}
	defer listenerFile.Close()

	// Create a net.Listener from the *os.File.
	ln, err := net.FileListener(listenerFile)
	if err != nil {
		return nil, err
	}

	return ln, nil
}

func createListener(addr string) (net.Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return ln, nil
}

func CreateOrImportListener(svcName, addr string) (net.Listener, error) {
	// Try and import a listener for addr. If it's found, use it.
	ln, err := importListener(svcName, addr)
	if err == nil {
		fmt.Printf("Imported listener file descriptor for %v.\n", addr)
		return ln, nil
	}

	// No listener was imported, that means this process has to create one.
	ln, err = createListener(addr)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Created listener file descriptor for %v.\n", addr)
	return ln, nil
}

// getListenerFile 获取监听的句柄
func getListenerFile(ln net.Listener) (*os.File, error) {
	switch t := ln.(type) {
	case *net.TCPListener:
		return t.File()
	case *net.UnixListener:
		return t.File()
	}
	return nil, fmt.Errorf("unsupported listener: %T", ln)
}

func forkChild(svcName, addr string, ln net.Listener) (*os.Process, error) {
	// Get the file descriptor for the listener and marshal the metadata to pass
	// to the child in the environment.
	lnFile, err := getListenerFile(ln)
	if err != nil {
		return nil, err
	}
	defer lnFile.Close()

	l := listener{
		Addr:     addr,
		FD:       3,
		Filename: lnFile.Name(),
	}
	listenerEnv, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}

	// Pass stdin, stdout, and stderr along with the listener to the child.
	files := []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
		lnFile,
	}

	// Get current environment and add in the listener to it.
	environment := append(os.Environ(), getSvcListenKey(svcName, addr)+string(listenerEnv))

	// Get current process name and directory.
	execName, err := os.Executable()
	if err != nil {
		return nil, err
	}
	execDir := filepath.Dir(execName)

	// Spawn child process.
	p, err := os.StartProcess(execName, []string{execName}, &os.ProcAttr{
		Dir:   execDir,
		Env:   environment,
		Files: files,
		Sys:   &syscall.SysProcAttr{},
	})
	if err != nil {
		return nil, err
	}

	return p, nil
}

func WaitForSignals(svcName, addr string, ln net.Listener, server *http.Server) error {
	signalCh := make(chan os.Signal, 5)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	for {
		select {
		case s := <-signalCh:
			fmt.Printf("%v signal received.\n", s)
			switch s {
			case syscall.SIGHUP:
				// Fork a child process.
				p, err := forkChild(svcName, addr, ln)
				if err != nil {
					fmt.Printf("Unable to fork child: %v.\n", err)
					continue
				}
				fmt.Printf("Forked child %v.\n", p.Pid)
			default:
				// 其他的信号量默认直接退出
				break
			}
		}
		// 注意这里一定要补上，这里要退出外部的for循环的
		break
	}
	// 创建一个默认的context，让其超时后自动退出
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 优雅关闭服务
	return server.Shutdown(ctx)
}
