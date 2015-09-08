package main

import (
	"fmt"
	"github.com/nkozyra/api/memorycache"
)

func main() {
	parameters := make(map[string]string)
	parameters["page"] = "1"
	parameters["search"] = "nathan"
	err, c := memorycache.Evaluate("test", "NEW VALUE TO SET!", 60, parameters)
	fmt.Println(c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Returned value is", c.Age, "seconds old")
	fmt.Println(c.Contents)
}
