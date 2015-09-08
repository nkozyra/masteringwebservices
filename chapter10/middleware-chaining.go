package main

import (
	"fmt"
	"net/http"
)

func PrimaryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello")
}

func MiddlewareHandler(h http.HandlerFunc) http.HandlerFunc {
	fmt.Println("Middleware!")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func WrappedMiddleware(h http.Handler) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func middleware(ph http.HandlerFunc, middleHandlers ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	fmt.Println("hello?")
	var next http.HandlerFunc = ph
	for _, mw := range middleHandlers {
		fmt.Println("Um?")
		next = mw(ph)
	}
	return next
}

func main() {
	x := middleware(PrimaryHandler, MiddlewareHandler, MiddlewareHandler, MiddlewareHandler)
	http.HandleFunc("/test", x)
	http.HandleFunc("/alternative", WrappedMiddleware)
	http.ListenAndServe(":9000", nil)
}
