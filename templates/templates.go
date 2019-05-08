package templates

import (
	"github.com/chanzuckerberg/fogg/config/v1"
	"github.com/gobuffalo/packr"
)

type T struct {
	Account          packr.Box
	Components       map[v1.ComponentKind]packr.Box
	Env              packr.Box
	Global           packr.Box
	Module           packr.Box
	ModuleInvocation packr.Box
	Repo             packr.Box
	TravisCI         packr.Box
}

var Templates = &T{
	Account: packr.NewBox("account"),
	Components: map[v1.ComponentKind]packr.Box{
		v1.ComponentKindTerraform:    packr.NewBox("component/terraform"),
		v1.ComponentKindHelmTemplate: packr.NewBox("component/helm_template"),
	},
	Env:              packr.NewBox("env"),
	Global:           packr.NewBox("global"),
	Module:           packr.NewBox("module"),
	ModuleInvocation: packr.NewBox("module-invocation"),
	Repo:             packr.NewBox("repo"),
	TravisCI:         packr.NewBox("travis-ci"),
}
