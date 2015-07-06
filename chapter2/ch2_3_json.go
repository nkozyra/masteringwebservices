package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type User struct {
	Name  string `xml:"name"`
	Email string `xml:"email"`
	ID    int    `xml:"id"`
}

func userRouter(w http.ResponseWriter, r *http.Request) {
	ourUser := User{}
	ourUser.Name = "Bill Smith"
	ourUser.Email = "bill.smith@example.com"
	ourUser.ID = 100

	output, _ := xml.Marshal(&ourUser)
	fmt.Fprintln(w, string(output))
}

func main() {

	fmt.Println("Starting JSON server")
	http.HandleFunc("/user", userRouter)
	http.ListenAndServe(":8080", nil)

}
