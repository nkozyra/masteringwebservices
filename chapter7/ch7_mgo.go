package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

func main() {
	session, err := mgo.Dial("localhost:2701")
	if err != nil {
		fmt.Println(err.Error())
	}

	c := session.DB("social-network").C("sessions")
	fmt.Println("hey")
}
