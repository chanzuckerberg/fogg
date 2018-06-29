package apply

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chanzuckerberg/fogg/plan"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/gobuffalo/packr"
	"github.com/hashicorp/hcl/hcl/printer"
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

	e = applyModules(fs, p.Modules, &tmp.Module)
	if e != nil {
		return e
	}

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

func applyModules(fs afero.Fs, p map[string]plan.Module, moduleBox *packr.Box) error {
	var e error
	for module, modulePlan := range p {
		path := fmt.Sprintf("%s/modules/%s", rootPath, module)
		e = fs.MkdirAll(path, 0755)
		if e != nil {
			return e
		}
		e = applyTree(moduleBox, afero.NewBasePathFs(fs, path), modulePlan)
		if e != nil {
			return e
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
	}
	return nil
}

func applyTree(source *packr.Box, dest afero.Fs, subst interface{}) (e error) {
	return source.Walk(func(path string, sourceFile packr.File) error {
		extension := filepath.Ext(path)
		basename := removeExtension(path)
		destExtension := filepath.Ext(basename)
		if extension == ".tmpl" {

			err := applyTemplate(sourceFile, dest, basename, subst)
			if err != nil {
				return errors.Wrap(err, "unable to apply template")
			}

		} else if extension == ".touch" {
			err := touchFile(dest, basename)
			if err != nil {
				return err
			}

		} else if extension == ".create" {
			err := createFile(dest, basename, sourceFile)
			if err != nil {
				return err
			}
		} else {
			log.Printf("copying %s", path)
			e = afero.WriteReader(dest, path, sourceFile)
			if e != nil {
				return errors.Wrap(e, "unable to copy file")
			}
		}

		if destExtension == ".tf" {
			fmtHcl(dest, basename)
		}

		return nil
	})

}

func fmtHcl(fs afero.Fs, path string) error {
	in, e := afero.ReadFile(fs, path)
	if e != nil {
		return e
	}
	out, e := printer.Format(in)
	if e != nil {
		return e
	}
	return afero.WriteReader(fs, path, bytes.NewReader(out))
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
