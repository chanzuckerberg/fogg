package apply

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/chanzuckerberg/fogg/config"
	"github.com/chanzuckerberg/fogg/plan"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/gobuffalo/packr"
	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const rootPath = "terraform"

func Apply(fs afero.Fs, conf *config.Config, tmp *templates.T, siccMode bool) error {
	p, err := plan.Eval(conf, siccMode, false)
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
		return errors.Wrap(e, "unable to apply global")
	}

	e = applyModules(fs, p.Modules, &tmp.Module)

	return errors.Wrap(e, "unable to apply modules")
}

func applyRepo(fs afero.Fs, p *plan.Plan, repoTemplates *packr.Box) error {
	e := applyTree(repoTemplates, fs, "", p.SiccMode, p)
	if e != nil {
		return e
	}
	// Remove after migration.

	if !p.SiccMode {
		e := fs.Remove(".sicc-version")
		if e != nil {
			log.Debug("error removing .sicc-version. ignoring")
		}
	}
	return nil
}

func applyGlobal(fs afero.Fs, p plan.Component, repoBox *packr.Box) error {
	path := fmt.Sprintf("%s/global", rootPath)
	e := fs.MkdirAll(path, 0755)
	if e != nil {
		return errors.Wrapf(e, "unable to make directory %s", path)
	}
	return applyTree(repoBox, fs, path, p.SiccMode, p)
}

func applyAccounts(fs afero.Fs, p *plan.Plan, accountBox *packr.Box) (e error) {
	for account, accountPlan := range p.Accounts {
		path := fmt.Sprintf("%s/accounts/%s", rootPath, account)
		e = fs.MkdirAll(path, 0755)
		if e != nil {
			return errors.Wrap(e, "unable to make directories for accounts")
		}
		e = applyTree(accountBox, fs, path, accountPlan.SiccMode, accountPlan)
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
			return errors.Wrapf(e, "unable to make path %s", path)
		}
		e = applyTree(moduleBox, fs, path, modulePlan.SiccMode, modulePlan)
		if e != nil {
			return errors.Wrap(e, "unable to apply tree")
		}
	}
	return nil
}

func applyEnvs(fs afero.Fs, p *plan.Plan, envBox *packr.Box, componentBox *packr.Box) (e error) {
	for env, envPlan := range p.Envs {
		path := fmt.Sprintf("%s/envs/%s", rootPath, env)
		e = fs.MkdirAll(path, 0755)
		if e != nil {
			return errors.Wrapf(e, "unable to make directory %s", path)
		}
		e := applyTree(envBox, fs, path, envPlan.SiccMode, envPlan)
		if e != nil {
			return errors.Wrap(e, "unable to apply templates to env")
		}
		for component, componentPlan := range envPlan.Components {
			path = fmt.Sprintf("%s/envs/%s/%s", rootPath, env, component)
			e = fs.MkdirAll(path, 0755)
			if e != nil {
				return errors.Wrap(e, "unable to make directories for component")
			}
			e := applyTree(componentBox, fs, path, componentPlan.SiccMode, componentPlan)
			if e != nil {
				return errors.Wrap(e, "unable to apply templates for component")
			}

			if componentPlan.BootstrapModule != "" {
				e := applyModule(fs, path, componentPlan.BootstrapModule, templates.Templates.ModuleInvocation)
				if e != nil {
					return errors.Wrap(e, "unable to apply module")
				}
			}

		}
	}
	return nil
}

func applyTree(source *packr.Box, dest afero.Fs, targetBasePath string, siccMode bool, subst interface{}) (e error) {
	return source.Walk(func(path string, sourceFile packr.File) error {
		extension := filepath.Ext(path)
		target := getTargetPath(targetBasePath, path, siccMode)

		targetExtension := filepath.Ext(target)
		if extension == ".tmpl" {
			e = applyTemplate(sourceFile, dest, target, subst)
			if e != nil {
				return errors.Wrap(e, "unable to apply template")
			}
		} else if extension == ".touch" {
			e = touchFile(dest, target)
			if e != nil {
				return errors.Wrapf(e, "unable to touch file %s", target)
			}
		} else if extension == ".create" {
			e = createFile(dest, target, sourceFile)
			if e != nil {
				return errors.Wrapf(e, "unable to create file %s", target)
			}
		} else {
			log.Infof("%s copied", path)
			e = afero.WriteReader(dest, path, sourceFile)
			if e != nil {
				return errors.Wrap(e, "unable to copy file")
			}
		}

		if !siccMode && target == "fogg.tf" {
			r := filepath.Join(targetBasePath, "sicc.tf")
			log.Infof("%s removed", r)
			e = dest.Remove(r)
			if e != nil {
				log.Debugf("error removing %s. ignoring", r)
			}
		}

		if targetExtension == ".tf" {
			e = fmtHcl(dest, target)
			if e != nil {
				return errors.Wrap(e, "unable to format HCL")
			}
		}
		// Some output files need to be executable, but we lose all file mode info
		// when the files get put in a Box. This is a blunt instrument to make
		// sure those work
		return dest.Chmod(target, 0644)
	})
}

func fmtHcl(fs afero.Fs, path string) error {
	in, e := afero.ReadFile(fs, path)
	if e != nil {
		return errors.Wrapf(e, "unable to read file %s", path)
	}
	out, e := printer.Format(in)
	if e != nil {
		return errors.Wrapf(e, "fmt hcl failed for %s", path)
	}
	return afero.WriteReader(fs, path, bytes.NewReader(out))
}

func touchFile(dest afero.Fs, path string) error {
	_, err := dest.Stat(path)
	if err != nil { // TODO we might not want to do this for all errors
		log.Infof("%s touched", path)
		_, err = dest.Create(path)
		if err != nil {
			return errors.Wrap(err, "unable to touch file")
		}
	} else {
		log.Infof("%s skipped touch", path)
	}
	return nil
}

func createFile(dest afero.Fs, path string, sourceFile io.Reader) error {
	_, err := dest.Stat(path)
	if err != nil { // TODO we might not want to do this for all errors
		log.Infof("%s created", path)
		err = afero.WriteReader(dest, path, sourceFile)
		if err != nil {
			return errors.Wrap(err, "unable to create file")
		}
	} else {
		log.Infof("%s skipped", path)
	}
	return nil
}

func removeExtension(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}

func applyTemplate(sourceFile io.Reader, dest afero.Fs, path string, overrides interface{}) error {
	log.Infof("%s templated", path)
	writer, err := dest.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return errors.Wrap(err, "unable to open file")
	}
	t := util.OpenTemplate(sourceFile)
	return t.Execute(writer, overrides)
}

// This should really be part of the plan stage, not apply. But going to
// leave it here for now and re-think it when we make this mechanism
// general purpose.
type moduleData struct {
	ModuleName   string
	ModuleSource string
	Variables    []string
	Outputs      []string
}

func applyModule(fs afero.Fs, path, mod string, box packr.Box) error {
	e := fs.MkdirAll(path, 0755)
	if e != nil {
		return errors.Wrapf(e, "couldn't create %s directory", path)
	}

	c, e := util.DownloadAndParseModule(mod)
	if e != nil {
		return errors.Wrap(e, "could not download or parse module")
	}

	// This should really be part of the plan stage, not apply. But going to
	// leave it here for now and re-think it when we make this mechanism
	// general purpose.
	variables := make([]string, 0)
	for _, v := range c.Variables {
		variables = append(variables, v.Name)
	}
	sort.Strings(variables)
	outputs := make([]string, 0)
	for _, o := range c.Outputs {
		outputs = append(outputs, o.Name)
	}
	sort.Strings(outputs)
	moduleName := filepath.Base(mod)
	re := regexp.MustCompile("\\?ref=.*")
	moduleName = re.ReplaceAllString(moduleName, "")

	f, e := box.Open("main.tf.tmpl")
	if e != nil {
		return errors.Wrap(e, "could not open template file")
	}

	e = applyTemplate(f, fs, filepath.Join(path, "main.tf"), &moduleData{moduleName, mod, variables, outputs})
	if e != nil {
		return errors.Wrap(e, "unable to apply template for main.tf")
	}
	e = fmtHcl(fs, filepath.Join(path, "main.tf"))
	if e != nil {
		return errors.Wrap(e, "unable to format main.tf")
	}

	f, e = box.Open("outputs.tf.tmpl")
	if e != nil {
		return errors.Wrap(e, "could not open template file")
	}

	e = applyTemplate(f, fs, filepath.Join(path, "outputs.tf"), &moduleData{moduleName, mod, variables, outputs})
	if e != nil {
		return errors.Wrap(e, "unable to apply template for outputs.tf")
	}

	e = fmtHcl(fs, filepath.Join(path, "outputs.tf"))
	if e != nil {
		return errors.Wrap(e, "unable to format outputs.tf")
	}

	return nil
}

func getTargetPath(basePath, path string, siccMode bool) string {
	target := filepath.Join(basePath, path)
	extension := filepath.Ext(path)

	if extension == ".tmpl" || extension == ".touch" || extension == ".create" {
		target = removeExtension(target)
	}
	if siccMode && filepath.Base(target) == "fogg.tf" {
		target = filepath.Join(filepath.Dir(target), "sicc.tf")
	}

	return target
}
