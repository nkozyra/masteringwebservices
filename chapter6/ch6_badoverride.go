package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

const (
	URL = "https://localhost/api/users"
)

func main() {

	customTransport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	customClient := &http.Client{Transport: customTransport}
	response, err := customClient.Get(URL)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(response)
	}

}
