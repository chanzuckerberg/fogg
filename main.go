package main

import (
	"flag"
	"os"

	finit "github.com/ryanking/fogg/init" // cannot import as init
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
		finit.Init(fs)
	case "plan":
		p, _ := plan.Plan(fs)
		plan.Print(p)
	}
	return
}
