package main

import (
	"flag"
	"os"

	"github.com/ryanking/fogg/plan"
	"github.com/spf13/afero"
)

func main() {
	flag.Parse()
	cmd := flag.Arg(0)

	pwd, _ := os.Getwd()
	fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

	switch cmd {
	case "init":
		Init(fs)
	case "plan":
		plan.Plan(fs)
	}
	return
}
