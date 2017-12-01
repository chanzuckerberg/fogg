package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/spf13/afero"
)

func main() {
	flag.Parse()
	cmd := flag.Arg(0)

	pwd, _ := os.Getwd()
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

	if cmd == "init" {
		Init(fs)
		return
	}
	f, _ := os.Open("fogg.json")
	c, _ := ReadConfig(io.ReadCloser(f))
	fmt.Println("hello world")
	fmt.Printf("%#v\n", c)
}
