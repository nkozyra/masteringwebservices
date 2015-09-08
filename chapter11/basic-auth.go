package main

import (
	"fmt"
	"github.com/gorilla/mux"
	//"log"
	"net/http"
)

func TestInterface(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	fmt.Fprintln(w, "hi")
}

func main() {

	Routes := mux.NewRouter()
	Routes.HandleFunc("/test", TestInterface).Methods("GET")
	http.ListenAndServe(":9000", Routes)

}
