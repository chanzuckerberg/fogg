package templates

import (
	"embed"
	"io/fs"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/sirupsen/logrus"
)

// NOTE(el): due to a design decision of go embed, we enumerate files starting with
//
//	the following set of characters: {.} giving explicit attention to directories with ONLY
//	files starting with said characters
//
//go:embed templates/.github/*
//go:embed templates/atlantis/*
//go:embed templates/pre-commit/*
//go:embed templates/pre-commit/root/.pre-commit-config.yaml.tmpl
//go:embed templates/circleci/.circleci/*
//go:embed templates/common/*
//go:embed templates/component/cdktf/.eslintrc.json.tmpl
//go:embed templates/component/cdktf/.prettierrc.json.tmpl
//go:embed templates/component/cdktf/.gitignore.create
//go:embed templates/component/cdktf/.gitattributes.create
//go:embed templates/component/*
//go:embed templates/env/*
//go:embed templates/module/cdktf/.eslintrc.json.tmpl
//go:embed templates/module/cdktf/.prettierrc.json.tmpl
//go:embed templates/module/cdktf/.gitignore.create
//go:embed templates/module/cdktf/.gitattributes.create
//go:embed templates/module/*
//go:embed templates/module/terraform/.update-readme.sh.rm
//go:embed templates/module-invocation/*
//go:embed templates/repo
//go:embed templates/repo/scripts/*
//go:embed templates/repo/.fogg-version.tmpl
//go:embed templates/repo/.gitattributes.tmpl
//go:embed templates/repo/.gitignore.tmpl
//go:embed templates/repo/.terraformignore.tmpl
//go:embed templates/repo/terraform.d/.keep.touch
//go:embed templates/repo/.terraform.d/plugin-cache/.gitignore
//go:embed templates/travis-ci/.travis.yml.tmpl
//go:embed templates/tfe/*
//go:embed templates/turbo/*
//go:embed templates/turbo/root/.vscode/settings.json.tmpl
var templates embed.FS

type T struct {
	Common           fs.FS
	Components       map[v2.ComponentKind]fs.FS
	Env              fs.FS
	Module           map[v2.ModuleKind]fs.FS
	ModuleInvocation fs.FS
	Repo             fs.FS
	TravisCI         fs.FS
	CircleCI         fs.FS
	GitHubActionsCI  fs.FS
	Atlantis         fs.FS
	PreCommitRoot    fs.FS
	PreCommitActions fs.FS
	TurboRoot        fs.FS
	TFE              fs.FS
}

// we control the inputs so should never panic
func mustFSSub(dir string) fs.FS {
	fs, err := fs.Sub(templates, dir)
	if err != nil {
		logrus.Fatalf("could not find templates for %s: %s", dir, err)
		return nil
	}
	return fs
}

var Templates = &T{
	Common: mustFSSub("templates/common"),
	Components: map[v2.ComponentKind]fs.FS{
		v2.ComponentKindTerraform: mustFSSub("templates/component/terraform"),
		v2.ComponentKindCDKTF:     mustFSSub("templates/component/cdktf"),
		v2.ComponentKindEnvtio:    mustFSSub("templates/component/envtio"),
	},
	Env: mustFSSub("templates/env"),
	Module: map[v2.ModuleKind]fs.FS{
		v2.ModuleKindTerraform: mustFSSub("templates/module/terraform"),
		v2.ModuleKindCDKTF:     mustFSSub("templates/module/cdktf"),
	},
	ModuleInvocation: mustFSSub("templates/module-invocation"),
	Repo:             mustFSSub("templates/repo"),
	TravisCI:         mustFSSub("templates/travis-ci"),
	CircleCI:         mustFSSub("templates/circleci"),
	GitHubActionsCI:  mustFSSub("templates/.github"),
	Atlantis:         mustFSSub("templates/atlantis"),
	PreCommitRoot:    mustFSSub("templates/pre-commit/root"),
	PreCommitActions: mustFSSub("templates/pre-commit/actions"),
	TurboRoot:        mustFSSub("templates/turbo/root"),
	TFE:              mustFSSub("templates/tfe"),
}
