package templates

import (
	"github.com/chanzuckerberg/fogg/config"
	"github.com/gobuffalo/packr"
)

type T struct {
	Account          packr.Box
	Component        map[config.ComponentKind]packr.Box
	Env              packr.Box
	Global           packr.Box
	Module           packr.Box
	ModuleInvocation packr.Box
	Repo             packr.Box
	TravisCI         packr.Box
}

var Templates = &T{
	Account: packr.NewBox("account"),
	Component: map[config.ComponentKind]packr.Box{
		config.ComponentKindTerraform:    packr.NewBox("component/terraform"),
		config.ComponentKindHelmTemplate: packr.NewBox("component/helm_template"),
	},
	Env:              packr.NewBox("env"),
	Global:           packr.NewBox("global"),
	Module:           packr.NewBox("module"),
	ModuleInvocation: packr.NewBox("module-invocation"),
	Repo:             packr.NewBox("repo"),
	TravisCI:         packr.NewBox("travis-ci"),
}
