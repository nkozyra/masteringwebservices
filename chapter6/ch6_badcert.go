package main

import (
	"fmt"
	"net/http"
)

const (
	URL = "https://localhost/api/users"
)

func main() {

	_, err := http.Get(URL)
	if err != nil {

		fmt.Println(err.Error())
	}

}
