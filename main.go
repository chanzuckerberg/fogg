package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	f, _ := os.Open("fogg.json")
	c, _ := ReadConfig(io.ReadCloser(f))
	fmt.Println("hello world")
	fmt.Printf("%#v\n", c)
}
