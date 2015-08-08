package main

import (
	"github.com/braintree/manners"
	"net/http"
)

func SignalListener() {

}

func main() {

	go func() {
		SignalListener()
	}()

	handler := MyHTTPHandler()
	server := manners.NewServer()
	server.ListenAndServe(":7000", handler)

}
