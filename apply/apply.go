package apply

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chanzuckerberg/fogg/plan"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/gobuffalo/packr"
	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/config/module"
	"github.com/hashicorp/terraform/svchost/auth"
	"github.com/hashicorp/terraform/svchost/disco"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const rootPath = "terraform"

func Apply(fs afero.Fs, configFile string, tmp *templates.T) error {
	p, err := plan.Eval(fs, configFile)
	if err != nil {
		return errors.Wrap(err, "unable to evaluate plan")
	}

	e := applyRepo(fs, p, &tmp.Repo)
	if e != nil {
		return errors.Wrap(e, "unable to apply repo")
	}

	e = applyAccounts(fs, p, &tmp.Account)
	if e != nil {
		return errors.Wrap(e, "unable to apply accounts")
	}

	e = applyEnvs(fs, p, &tmp.Env, &tmp.Component)
	if e != nil {
		return errors.Wrap(e, "unable to apply envs")
	}

	e = applyGlobal(fs, p.Global, &tmp.Global)
	if e != nil {
		return e
	}

	// TODO modules

	return nil
}

func applyRepo(fs afero.Fs, p *plan.Plan, repoBox *packr.Box) error {
	return applyTree(repoBox, fs, p)
}

func applyGlobal(fs afero.Fs, p plan.Component, repoBox *packr.Box) error {
	path := fmt.Sprintf("%s/global", rootPath)
	e := fs.MkdirAll(path, 0755)
	if e != nil {
		return e
	}
	return applyTree(repoBox, afero.NewBasePathFs(fs, path), p)
}

func applyAccounts(fs afero.Fs, p *plan.Plan, accountBox *packr.Box) (e error) {
	for account, accountPlan := range p.Accounts {
		path := fmt.Sprintf("%s/accounts/%s", rootPath, account)
		e = fs.MkdirAll(path, 0755)
		if e != nil {
			return errors.Wrap(e, "unable to make directories for accounts")
		}
		e = applyTree(accountBox, afero.NewBasePathFs(fs, path), accountPlan)
		if e != nil {
			return errors.Wrap(e, "unable to apply templates to account")
		}
	}
	return nil
}

func applyEnvs(fs afero.Fs, p *plan.Plan, envBox *packr.Box, componentBox *packr.Box) (e error) {
	for env, envPlan := range p.Envs {
		path := fmt.Sprintf("%s/envs/%s", rootPath, env)
		e = fs.MkdirAll(path, 0755)
		if e != nil {
			return errors.Wrap(e, "unable to make directies for envs")
		}
		e := applyTree(envBox, afero.NewBasePathFs(fs, path), envPlan)
		if e != nil {
			return errors.Wrap(e, "unable to apply templates to env")
		}
		for component, componentPlan := range envPlan.Components {
			path := fmt.Sprintf("%s/envs/%s/%s", rootPath, env, component)
			e = fs.MkdirAll(path, 0755)
			if e != nil {
				return errors.Wrap(e, "unable to make directories for component")
			}
			e := applyTree(componentBox, afero.NewBasePathFs(fs, path), componentPlan)
			if e != nil {
				return errors.Wrap(e, "unable to apply templates for component")
			}
		}
		path = filepath.Join(rootPath, "envs", env, "cloud-env")
		if envPlan.Type == "aws" {
			applyModule(fs, path, "git@github.com:chanzuckerberg/shared-infra//terraform/modules/aws-env")
		}
	}
	return nil
}

func applyTree(source *packr.Box, dest afero.Fs, subst interface{}) (e error) {
	return source.Walk(func(path string, sourceFile packr.File) error {
		extension := filepath.Ext(path)
		basename := removeExtension(path)
		if extension == ".tmpl" {

			err := applyTemplate(sourceFile, dest, basename, subst)
			if err != nil {
				return errors.Wrap(err, "unable to apply template")
			}

			//     if dest.endswith('.tf'):
			//         subprocess.call(['terraform', 'fmt', dest])
		} else if extension == ".touch" {
			touchFile(dest, basename)
			//     if dest.endswith('.tf'):
			//         subprocess.call(['terraform', 'fmt', dest])

		} else if extension == ".create" {
			createFile(dest, basename, sourceFile)
			//     if dest.endswith('.tf'):
			//         subprocess.call(['terraform', 'fmt', dest])

		} else {
			log.Printf("copying %s", path)
			e = afero.WriteReader(dest, path, sourceFile)
			if e != nil {
				return errors.Wrap(e, "unable to copy file")
			}
		}
		return nil
	})

}

func touchFile(dest afero.Fs, path string) error {
	_, err := dest.Stat(path)
	if err != nil { // TODO we might not want to do this for all errors
		log.Printf("touching %s", path)
		_, err = dest.Create(path)
		if err != nil {
			return errors.Wrap(err, "unable to touch file")
		}
	} else {
		log.Printf("skipping touch on existing file %s", path)
	}
	return nil
}

func createFile(dest afero.Fs, path string, sourceFile io.Reader) error {
	_, err := dest.Stat(path)
	if err != nil { // TODO we might not want to do this for all errors
		log.Printf("creating %s", path)
		err = afero.WriteReader(dest, path, sourceFile)
		if err != nil {
			return errors.Wrap(err, "unable to create file")
		}
	} else {
		log.Printf("skipping create on existing file %s", path)
	}
	return nil
}

func removeExtension(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}

func applyTemplate(sourceFile io.Reader, dest afero.Fs, path string, overrides interface{}) error {
	log.Printf("templating %s", path)
	writer, err := dest.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return errors.Wrap(err, "unable to open file")
	}
	t := util.OpenTemplate(sourceFile)
	return t.Execute(writer, overrides)
}

func downloadModule(dir, mod string) error {
	disco := disco.NewDisco()
	s := module.NewStorage(dir, disco, auth.NoCredentials)

	return s.GetModule(dir, mod)
}

func downloadAndParseModule(dir, mod string) (*config.Config, error) {
	e := downloadModule(dir, mod)
	if e != nil {
		return nil, errors.Wrap(e, "unable to download module")
	}
	return config.LoadDir(dir)
}

type moduleData struct {
	ModuleName   string
	ModuleSource string
	Variables    []string
	Outputs      []string
}

func applyModule(fs afero.Fs, path, mod string) error {
	e := fs.MkdirAll(path, 0755)
	if e != nil {
		return errors.Wrapf(e, "couldn't create %s directory", path)
	}

	dir, e := ioutil.TempDir("", "fogg")
	if e != nil {
		return e
	}
	c, e := downloadAndParseModule(dir, mod)
	variables := make([]string, 0)
	for _, v := range c.Variables {
		variables = append(variables, v.Name)
	}
	outputs := make([]string, 0)
	for _, o := range c.Outputs {
		outputs = append(outputs, o.Name)
	}
	moduleName := filepath.Base(mod)

	main := `
module "{{.ModuleName}}" {
  source = "{{.ModuleSource}}"
  {{range .Variables -}}
    {{.}} = "${var.{{.}}}"
  {{ end}}
}
`
	mainTemp := template.Must(template.New("tmpl").Parse(string(main)))
	mainFile, e := fs.OpenFile(filepath.Join(path, "main.tf"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if e != nil {
		return errors.Wrapf(e, "unable to open %s", filepath.Join(path, "main.tf"))
	}
	mainTemp.Execute(mainFile, &moduleData{moduleName, mod, variables, outputs})

	output := `
{{ $outer := . -}}
{{- range .Outputs -}}
output "{{.}}" {
  value = "${module.{{$outer.ModuleName}}.{{.}}}"
}

{{end}}`
	outputsTemp := template.Must(template.New("tmpl").Parse(string(output)))
	outputFile, e := fs.OpenFile(filepath.Join(path, "outputs.tf"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if e != nil {
		return errors.Wrapf(e, "unable to open %s", filepath.Join(path, "outputs.tf"))
	}

	outputsTemp.Execute(outputFile, &moduleData{moduleName, mod, variables, outputs})

	return e
}
