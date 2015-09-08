package main

import (
	"fmt"
	"github.com/nkozyra/api/ratelimit"
)

func main() {
	valid := ratelimit.CheckRequest("127.0.0.1")
	fmt.Println(valid)
}
