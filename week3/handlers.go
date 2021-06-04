package main

import (
	"net/http"
	"fmt"
	"time"
)
// Hello handleFunc
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s", r.Host)
}

// Bye handlerFunc
func ByeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bye, %s", r.Host)
}

// DoSomething handlerFunc
func DoSomethingHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	fmt.Fprintf(w, "Done!")
}