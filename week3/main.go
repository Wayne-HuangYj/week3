package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	_ "net/http/pprof"
	"context"
	"errors"
	"time"

	"golang.org/x/sync/errgroup"
)

// 根据自定义Server开启http服务
func startServe(ctx context.Context, addr string, handler http.Handler) error {
	server := http.Server {
		Addr: addr,
		Handler: handler,
	}

	// 模拟报错
	if _, ok := handler.(errHandler); ok {
		time.Sleep(5 * time.Second)
		return errors.New("Error occurred")
	}

	go func() {
		select {
		case <-ctx.Done():
			log.Println(ctx.Err())
		}
		server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}


func main() {
	// 创建监听kill信号的管道
	signalStop := make(chan os.Signal, 1)
	signal.Notify(signalStop)
	// 创建一个用于防止goroutine泄露的context
	ctx, cancel := context.WithCancel(context.Background())
	// 创建errgroup
	g, ctx := errgroup.WithContext(ctx)
	// 开启第一个Hello服务
	g.Go(func() error {
		return startServe(ctx, ":8000", helloHandler{})
	})
	// 开启第二个Bye服务
	g.Go(func() error {
		return startServe(ctx, ":8002", byeHandler{})
	})
	// 开启第三个error服务
	g.Go(func() error {
		return startServe(ctx, ":8004", errHandler{})
	})

	// 监听signal的，如果收到信号，应该将context关闭以关闭所有goroutine
	g.Go(func() error {
		select {
		case signal := <-signalStop:
			// log.Printf("A %s signal was received",  signal)
			return fmt.Errorf("A %s signal was received",  signal)
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	if err := g.Wait(); err != nil { 
		// 有任何一个goroutine出错了，或者kill信号监听，直接关闭ctx来关闭所有goroutine
		cancel()
		log.Println(err)
	}
}