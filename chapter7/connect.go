package main

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	mc := memcache.New("[host]:11211")
	mc.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})
	it, err := mc.Get("foo")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(it)
	fmt.Println("hello!")
}
