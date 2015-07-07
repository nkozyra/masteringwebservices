package main

import (
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "interface-update.html")
	})
	http.ListenAndServe(":9000", nil)

}
