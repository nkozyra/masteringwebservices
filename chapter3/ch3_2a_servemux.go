package main

import (
	"fmt"
	"net/http"
	"regexp"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path
		message := "You have triggered nothing"

		testMatch, _ := regexp.MatchString("/testing[0-9]{3}", path)

		if testMatch == true {
			// helper functions
			message = "You hit the test!"
		}

		fmt.Fprintln(w, message)
	})

	http.ListenAndServe(":8080", nil)
}
