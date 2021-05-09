package main

import (
	"net/http"
	"fmt"
)
// 其中一个handler
type helloHandler struct{

}

func (h helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s", r.Host)
}

// 另外一个handler
type byeHandler struct {

}

func (h byeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bye, %s", r.Host)
}

// 最后一个handler模拟http服务报错
type errHandler struct {

}

func (h errHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Error occurred!")
}