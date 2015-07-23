package fmttest

import (
	"fmt"
)

type format int

func (f format) Test() {
	fmt.Println("YES")
}

const (
	Bold   format = 10
	Italic format = 20
	Strong format = 30
)

func CheckFormat() {
	fmt.Println(Bold)
}
