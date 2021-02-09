package templates

import (
	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/gobuffalo/packr"
)

type T struct {
	Common           packr.Box
	Components       map[v2.ComponentKind]packr.Box
	Env              packr.Box
	Module           packr.Box
	ModuleInvocation packr.Box
	Repo             packr.Box
	TravisCI         packr.Box
	CircleCI         packr.Box
	GitHubActionsCI  packr.Box
}

var Templates = &T{
	Common: packr.NewBox("common"),
	Components: map[v2.ComponentKind]packr.Box{
		v2.ComponentKindTerraform:    packr.NewBox("component/terraform"),
		v2.ComponentKindHelmTemplate: packr.NewBox("component/helm_template"),
	},
	Env:              packr.NewBox("env"),
	Module:           packr.NewBox("module"),
	ModuleInvocation: packr.NewBox("module-invocation"),
	Repo:             packr.NewBox("repo"),
	TravisCI:         packr.NewBox("travis-ci"),
	CircleCI:         packr.NewBox("circleci"),
	GitHubActionsCI:  packr.NewBox(".github"),
}
