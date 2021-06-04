package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	_ "net/http/pprof"
	"context"
	"golang.org/x/sync/errgroup"
)
// 多路复用处理函数，ServerMux也是一种Http的handler，它内部也实现了serveHttp的函数，用这个去代替defaultServerMux
// 它负责解析URL，根据URL pattren解析，把请求发放到不同的handler中处理，因此可以把一些handler和url的pattern注册到自己的mux中
type Server struct {
	mux *http.ServeMux
	server *http.Server
	addr string
	ctx context.Context
}

// 创建一个自定义的Server，传入一个context用作控制，并且指定监听的host
func NewServer(ctx context.Context, addr string) (server *Server) {
	server = &Server{
		mux: http.NewServeMux(),
		addr: addr,
		ctx: ctx,
	}
	// http.Server先不初始化，等到HandleFunc调用完了，用户要ListenAndServe的时候再初始化并且监听
	return
}

// 注册函数，将handler注册到自定义的ServeMux中
func (s *Server) HandleFunc(pattern string, handler func (http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
}

// 开启Server，监听服务
func (s *Server) ListenAndServe()  error {
	s.server = &http.Server {
		Addr: s.addr,
		Handler: s.mux,
	}

	go func() {
		select {
		case <-s.ctx.Done():  //  表示有请求发生的报错，通过context关闭所有的监听
			log.Println(s.ctx.Err())
		}
		s.server.Shutdown(context.Background())
	}()
	return s.server.ListenAndServe()
}

/*
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
*/


func main() {
	// 创建监听kill信号的管道
	signalStop := make(chan os.Signal, 1)
	signal.Notify(signalStop)
	// 创建一个用于防止goroutine泄露的context，如果发现有任何一个请求发生了报错，立刻停止ListenAndServer
	ctx, cancel := context.WithCancel(context.Background())
	// 创建errgroup
	g, ctx := errgroup.WithContext(ctx)

	// 创建一个server
	server := NewServer(ctx, ":8000")
	
	// 注册三个handler
	server.HandleFunc("/hello", HelloHandler)
	server.HandleFunc("/bye", ByeHandler)
	server.HandleFunc("/do", DoSomethingHandler)
	
	// 开启监听命令行signal的goroutine
	g.Go(func() error {
		select {
		case received := <-signalStop:
			return fmt.Errorf("A %s signal was received",  received)
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	// 再开启HTTP的ListenAndServe
	g.Go(server.ListenAndServe)

	// 监听errorgroup有没有错误返回
	// 错误有两种，一种是由于收到signal而返回的
	// 第二种是由于HTTP请求出错了，导致ListenAndServe返回了一个错误
	if err := g.Wait(); err != nil { 
		// 为什么要用errorgroup，因为它有一个context，只要把这个context传递给各个goroutine
		// 当其中一个出错的时候会被errorgroup捕获到
		// c此时只需要简单地调用context配套的cancel函数，就可以把所有goroutine关闭
		// 前提就是goroutine必须要对这个ctx.Done()管道进行监听
		cancel()
		log.Printf("HTTP服务结束，原因：%v\n", err)
	}
}

// 废弃代码
// g.Go(func() error {
	// 	return startServe(ctx, ":8000", helloHandler{})
	// })
	// // 开启第二个Bye服务
	// g.Go(func() error {
	// 	return startServe(ctx, ":8002", byeHandler{})
	// })
	// // 开启第三个error服务
	// g.Go(func() error {
	// 	return startServe(ctx, ":8004", errHandler{})
	// })