package templates

import (
	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/gobuffalo/packr/v2"
)

type T struct {
	Common           *packr.Box
	Components       map[v2.ComponentKind]*packr.Box
	Env              *packr.Box
	Module           *packr.Box
	ModuleInvocation *packr.Box
	Repo             *packr.Box
	TravisCI         *packr.Box
	CircleCI         *packr.Box
	GitHubActionsCI  *packr.Box
}

var Templates = &T{
	Common: packr.New("common", "./common"),
	Components: map[v2.ComponentKind]*packr.Box{
		v2.ComponentKindTerraform:    packr.New("component/terraform", "./component/terraform"),
		v2.ComponentKindHelmTemplate: packr.New("component/helm_template", "./component/helm_template"),
	},
	Env:              packr.New("env", "./env"),
	Module:           packr.New("module", "./module"),
	ModuleInvocation: packr.New("module-invocation", "./module-invocation"),
	Repo:             packr.New("repo", "./repo"),
	TravisCI:         packr.New("travis-ci", "./travis-ci"),
	CircleCI:         packr.New("circleci", "./circleci"),
	GitHubActionsCI:  packr.New(".github", "./.github"),
}
