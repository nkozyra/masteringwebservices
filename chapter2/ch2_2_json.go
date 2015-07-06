package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	ID    int    `json:"int"`
}

func userRouter(w http.ResponseWriter, r *http.Request) {
	ourUser := User{}
	ourUser.Name = "Bill Smith"
	ourUser.Email = "bill.smith@example.com"
	ourUser.ID = 100

	output, _ := json.Marshal(&ourUser)
	fmt.Fprintln(w, string(output))
}

func main() {

	fmt.Println("Starting JSON server")
	http.HandleFunc("/user", userRouter)
	http.ListenAndServe(":8080", nil)

}
