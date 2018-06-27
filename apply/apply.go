package apply

import (
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
	"github.com/spf13/afero"
)

const rootPath = "terraform"

func Apply(fs afero.Fs, configFile string, tmp *templates.T) error {
	p, err := plan.Eval(fs, configFile)
	if err != nil {
		return err
	}

	e := applyRepo(fs, p, &tmp.Repo)
	if e != nil {
		return e
	}

	e = applyAccounts(fs, p, &tmp.Account)
	if e != nil {
		return e
	}

	e = applyEnvs(fs, p, &tmp.Env, &tmp.Component)
	if e != nil {
		return e
	}

	e = applyGlobal(fs, p.Global, &tmp.Global)

	// TODO modules

	return nil
}

func applyRepo(fs afero.Fs, p *plan.Plan, repoBox *packr.Box) error {
	return applyTree(repoBox, fs, p)
}

func applyGlobal(fs afero.Fs, p plan.Component, repoBox *packr.Box) error {
	log.Println("global")
	util.Dump(p)
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
			return e
		}
		e = applyTree(accountBox, afero.NewBasePathFs(fs, path), accountPlan)
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
			return e
		}
		e := applyTree(envBox, afero.NewBasePathFs(fs, path), envPlan)
		if e != nil {
			return e
		}
		for component, componentPlan := range envPlan.Components {
			path := fmt.Sprintf("%s/envs/%s/%s", rootPath, env, component)
			e = fs.MkdirAll(path, 0755)
			if e != nil {
				return e
			}
			e := applyTree(componentBox, afero.NewBasePathFs(fs, path), componentPlan)
			if e != nil {
				return e
			}
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
				return err
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
				return e
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
			return err
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
			return err
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
		return err
	}
	t := util.OpenTemplate(sourceFile)
	return t.Execute(writer, overrides)
}
