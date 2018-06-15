package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	finit "github.com/ryanking/fogg/init" // cannot import as init
	"github.com/ryanking/fogg/plan"
	"github.com/ryanking/fogg/util"
	"github.com/spf13/afero"
)

var (
	Version string
	GitSha  string
	Release string
	Dirty   string
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
	case "version":
		release, _ := strconv.ParseBool(Release)
		dirty, _ := strconv.ParseBool(Dirty)
		fmt.Println(util.VersionString(Version, GitSha, release, dirty))
	}
	return
}
