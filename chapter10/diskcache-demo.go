package main

import (
	"fmt"
	"github.com/nkozyra/api/diskcache"
)

func main() {
	err, c := diskcache.Evaluate("test", "Here is a value that will only live for 1 minute", 60)
	fmt.Println(c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Returned value is", c.Age, "seconds old")
	fmt.Println(c.Contents)
}
