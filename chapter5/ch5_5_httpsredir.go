package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

const serverName = "localhost"
const SSLport = ":443"
const HTTPport = ":8081"
const SSLprotocol = "https://"
const HTTPprotocol = "http://"

func secureRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You have arrived at port 443, and now you are now marginally more secure.")
}

func redirectNonSecure(w http.ResponseWriter, r *http.Request) {
	log.Println("Non-secure request initiated, redirecting.")
	redirectURL := SSLprotocol + serverName + r.RequestURI
	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}

func main() {

	wg := sync.WaitGroup{}

	log.Println("Starting redirection server, try to access @ http:")

	wg.Add(1)
	go func() {
		http.ListenAndServe(HTTPport, http.HandlerFunc(redirectNonSecure))
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		http.ListenAndServeTLS(SSLport, "cert.pem", "key.pem", http.HandlerFunc(secureRequest))
		//http.ListenAndServe(SSLport,http.HandlerFunc(secureRequest))
		wg.Done()
	}()

	wg.Wait()

}
