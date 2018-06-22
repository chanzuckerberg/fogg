package apply

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
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

	applyRepo(fs, p, &tmp.Repo)
	applyAccounts(fs, p, &tmp.Account)
	applyEnvs(fs, p, &tmp.Env)

	return nil
}

func applyRepo(fs afero.Fs, p *plan.Plan, repoBox *packr.Box) error {
	applyTree(repoBox, fs, p)
	return nil
}

func applyAccounts(fs afero.Fs, p *plan.Plan, accountBox *packr.Box) error {
	for account, accountPlan := range p.Accounts {
		path := fmt.Sprintf("%s/accounts/%s", rootPath, account)
		fs.MkdirAll(path, 0755)
		applyTree(accountBox, afero.NewBasePathFs(fs, path), accountPlan)
	}
	return nil
}

func applyEnvs(fs afero.Fs, p *plan.Plan, envBox *packr.Box) error {
	for env, envPlan := range p.Envs {
		path := fmt.Sprintf("%s/envs/%s", rootPath, env)
		fs.MkdirAll(path, 0755)
		applyTree(envBox, afero.NewBasePathFs(fs, path), envPlan)
	}
	return nil
}

func applyTree(source *packr.Box, dest afero.Fs, subst interface{}) error {
	source.Walk(func(path string, sourceFile packr.File) error {
		extension := filepath.Ext(path)
		if extension == ".tmpl" {

			err := applyTemplate(sourceFile, dest, path, subst)
			if err != nil {
				panic(err)
			}

			//     if dest.endswith('.tf'):
			//         subprocess.call(['terraform', 'fmt', dest])
		} else if extension == ".touch" {
			d := removeExtension(path)
			_, err := dest.Stat(d)
			if err != nil { // TODO we might not want to do this for all errors
				log.Printf("touching %s", d)
				dest.Create(d)
			} else {
				log.Printf("skipping touch on existing file %s", d)
			}
			//     if dest.endswith('.tf'):
			//         subprocess.call(['terraform', 'fmt', dest])

		} else if extension == ".create" {
			d := removeExtension(path)
			_, err := dest.Stat(d)
			if err != nil { // TODO we might not want to do this for all errors
				log.Printf("creating %s", d)
				afero.WriteReader(dest, path, sourceFile)
			} else {
				log.Printf("skipping create on existing file %s", d)
			}
			//     if dest.endswith('.tf'):
			//         subprocess.call(['terraform', 'fmt', dest])

		} else {
			log.Printf("copying %s", path)
			afero.WriteReader(dest, path, sourceFile)
		}
		return nil
	})
	return nil
}

func removeExtension(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}

func joinEnvs(m map[string]plan.Env) string {
	keys := make([]string, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, " ")
}

func applyTemplate(sourceFile packr.File, dest afero.Fs, path string, overrides interface{}) error {
	d := removeExtension(path)
	log.Printf("templating %s", d)
	writer, err := dest.OpenFile(d, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	t := util.OpenTemplate(sourceFile)
	return t.Execute(writer, overrides)
}
