package util

import "fmt"

func Dump(foo interface{}) {
	fmt.Printf("%#v\n", foo)
}
