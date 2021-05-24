package templates

import (
	"embed"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
)

// go:embed common
var common embed.FS

// go:embed component/terraform
var componentTerraform embed.FS

// go:embed env
var env embed.FS

// go:embed module
var module embed.FS

// go:embed module-invocation
var moduleInvokation embed.FS

// go:embed repo
var repo embed.FS

// go:embed travis-ci
var travisCI embed.FS

// go:embed circleci
var circleCI embed.FS

// go:embed .github
var gitHubActionsCI embed.FS

type T struct {
	Common           embed.FS
	Components       map[v2.ComponentKind]embed.FS
	Env              embed.FS
	Module           embed.FS
	ModuleInvocation embed.FS
	Repo             embed.FS
	TravisCI         embed.FS
	CircleCI         embed.FS
	GitHubActionsCI  embed.FS
}

var Templates = &T{
	Common: common,
	Components: map[v2.ComponentKind]embed.FS{
		v2.ComponentKindTerraform: componentTerraform,
	},
	Env:              env,
	Module:           module,
	ModuleInvocation: moduleInvokation,
	Repo:             repo,
	TravisCI:         travisCI,
	CircleCI:         circleCI,
	GitHubActionsCI:  gitHubActionsCI,
}
