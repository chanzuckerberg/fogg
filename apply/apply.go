package apply

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/plan"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/gobuffalo/packr"
	getter "github.com/hashicorp/go-getter"
	"github.com/hashicorp/hcl/hcl/printer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const rootPath = "terraform"

// Apply will run a plan and apply all the changes to the current repo.
func Apply(fs afero.Fs, conf *v2.Config, tmp *templates.T, upgrade bool) error {

	// TODO add ability to pass in pwd

	if !upgrade {
		toolVersion, err := util.VersionString()
		if err != nil {
			return err
		}
		versionChange, repoVersion, _ := checkToolVersions(fs, toolVersion)
		if versionChange {
			return errs.NewUserf("fogg version (%s) is different than version currently used to manage repo (%s). To upgrade add --upgrade.", toolVersion, repoVersion)
		}
	}
	p, err := plan.Eval(conf)
	if err != nil {
		return errs.WrapUser(err, "unable to evaluate plan")
	}

	e := applyRepo(fs, p, &tmp.Repo)
	if e != nil {
		return errs.WrapUser(e, "unable to apply repo")
	}

	if p.TravisCI.Enabled {
		e = applyTree(fs, &tmp.TravisCI, "", p.TravisCI)
		if e != nil {
			return errs.WrapUser(e, "unable to apply travis ci")
		}
	}

	e = applyAccounts(fs, p, &tmp.Account)
	if e != nil {
		return errs.WrapUser(e, "unable to apply accounts")
	}

	e = applyEnvs(fs, p, &tmp.Env, tmp.Components)
	if e != nil {
		return errs.WrapUser(e, "unable to apply envs")
	}

	e = applyGlobal(fs, p.Global, &tmp.Global)
	if e != nil {
		return errs.WrapUser(e, "unable to apply global")
	}

	e = applyModules(fs, p.Modules, &tmp.Module)
	return errs.WrapUser(e, "unable to apply modules")
}

func checkToolVersions(fs afero.Fs, current string) (bool, string, error) {
	f, e := fs.Open(".fogg-version")
	if e != nil {
		return false, "", errs.WrapUser(e, "unable to open .fogg-version file")
	}
	reader := io.ReadCloser(f)
	defer reader.Close()

	b, e := ioutil.ReadAll(reader)
	if e != nil {
		return false, "", errs.WrapUser(e, "unable to read .fogg-version file")
	}
	repoVersion := string(b)
	changed, e := versionIsChanged(repoVersion, current)
	return changed, repoVersion, e
}

func versionIsChanged(repo string, tool string) (bool, error) {
	repoVersion, repoSha, repoDirty := util.ParseVersion(repo)
	toolVersion, toolSha, toolDirty := util.ParseVersion(tool)

	if repoDirty || toolDirty {
		return true, nil
	}

	if (repoSha != "" || toolSha != "") && repoSha != toolSha {
		return true, nil
	}

	return toolVersion.NE(repoVersion), nil
}

func applyRepo(fs afero.Fs, p *plan.Plan, repoTemplates *packr.Box) error {
	return applyTree(fs, repoTemplates, "", p)
}

func applyGlobal(fs afero.Fs, p plan.Component, repoBox *packr.Box) error {
	log.Debug("applying global")
	path := fmt.Sprintf("%s/global", rootPath)
	e := fs.MkdirAll(path, 0755)
	if e != nil {
		return errs.WrapUserf(e, "unable to make directory %s", path)
	}
	return applyTree(fs, repoBox, path, p)
}

func applyAccounts(fs afero.Fs, p *plan.Plan, accountBox *packr.Box) (e error) {
	for account, accountPlan := range p.Accounts {
		path := fmt.Sprintf("%s/accounts/%s", rootPath, account)
		e = fs.MkdirAll(path, 0755)
		if e != nil {
			return errs.WrapUser(e, "unable to make directories for accounts")
		}
		e = applyTree(fs, accountBox, path, accountPlan)
		if e != nil {
			return errs.WrapUser(e, "unable to apply templates to account")
		}
	}
	return nil
}

func applyModules(fs afero.Fs, p map[string]plan.Module, moduleBox *packr.Box) (e error) {
	for module, modulePlan := range p {
		path := fmt.Sprintf("%s/modules/%s", rootPath, module)
		e = fs.MkdirAll(path, 0755)
		if e != nil {
			return errs.WrapUserf(e, "unable to make path %s", path)
		}
		e = applyTree(fs, moduleBox, path, modulePlan)
		if e != nil {
			return errs.WrapUser(e, "unable to apply tree")
		}
	}
	return nil
}

func applyEnvs(fs afero.Fs, p *plan.Plan, envBox *packr.Box, componentBoxes map[v1.ComponentKind]packr.Box) (e error) {
	log.Debug("applying envs")
	for env, envPlan := range p.Envs {
		log.Debugf("applying %s", env)
		path := fmt.Sprintf("%s/envs/%s", rootPath, env)
		e = fs.MkdirAll(path, 0755)
		if e != nil {
			return errs.WrapUserf(e, "unable to make directory %s", path)
		}
		e := applyTree(fs, envBox, path, envPlan)
		if e != nil {
			return errs.WrapUser(e, "unable to apply templates to env")
		}
		for component, componentPlan := range envPlan.Components {
			path = fmt.Sprintf("%s/envs/%s/%s", rootPath, env, component)
			e = fs.MkdirAll(path, 0755)
			if e != nil {
				return errs.WrapUser(e, "unable to make directories for component")
			}
			componentBox := componentBoxes[componentPlan.Kind.GetOrDefault()]
			e := applyTree(fs, &componentBox, path, componentPlan)
			if e != nil {
				return errs.WrapUser(e, "unable to apply templates for component")
			}

			if componentPlan.ModuleSource != nil {
				e := applyModuleInvocation(fs, path, *componentPlan.ModuleSource, templates.Templates.ModuleInvocation)
				if e != nil {
					return errs.WrapUser(e, "unable to apply module invocation")
				}
			}

		}
	}
	return nil
}

func applyTree(dest afero.Fs, source *packr.Box, targetBasePath string, subst interface{}) (e error) {
	return source.Walk(func(path string, sourceFile packr.File) error {

		extension := filepath.Ext(path)
		target := getTargetPath(targetBasePath, path)
		targetExtension := filepath.Ext(target)

		if extension == ".tmpl" {
			e = applyTemplate(sourceFile, dest, target, subst)
			if e != nil {
				return errs.WrapUser(e, "unable to apply template")
			}
		} else if extension == ".touch" {
			e = touchFile(dest, target)
			if e != nil {
				return errs.WrapUserf(e, "unable to touch file %s", target)
			}
		} else if extension == ".create" {
			e = createFile(dest, target, sourceFile)
			if e != nil {
				return errs.WrapUserf(e, "unable to create file %s", target)
			}
		} else if extension == ".rm" {
			e = os.Remove(target)
			if e != nil && !os.IsNotExist(e) {
				return errs.WrapUserf(e, "unable to remove %s", target)
			}
			log.Infof("%s removed", target)
		} else if extension == ".ln" {

			linkTargetBytes, err := ioutil.ReadAll(sourceFile)
			if err != nil {
				return errs.WrapUserf(err, "could not read source file %#v", sourceFile)
			}

			linkTarget := string(linkTargetBytes)

			err = linkFile(dest, target, linkTarget)

			if err != nil {
				return errs.WrapInternal(err, "can't symlink file")
			}

		} else {
			e = afero.WriteReader(dest, target, sourceFile)
			if e != nil {
				return errs.WrapUser(e, "unable to copy file")
			}
			log.Infof("%s copied", target)
		}

		if targetExtension == ".tf" {
			e = fmtHcl(dest, target)
			if e != nil {
				return errs.WrapUser(e, "unable to format HCL")
			}
		}
		return nil
	})
}

func fmtHcl(fs afero.Fs, path string) error {
	in, e := afero.ReadFile(fs, path)
	if e != nil {
		return errs.WrapUserf(e, "unable to read file %s", path)
	}
	out, e := printer.Format(in)
	if e != nil {
		return errs.WrapUserf(e, "fmt hcl failed for %s", path)
	}
	return afero.WriteReader(fs, path, bytes.NewReader(out))
}

func touchFile(dest afero.Fs, path string) error {
	dir, _ := filepath.Split(path)
	ospath := filepath.FromSlash(dir)
	err := dest.MkdirAll(ospath, 0755)
	if err != nil {
		return errs.WrapUserf(err, "couldn't create %s directory", dir)
	}

	_, err = dest.Stat(path)
	if err != nil { // TODO we might not want to do this for all errors
		log.Infof("%s touched", path)
		_, err = dest.Create(path)
		if err != nil {
			return errs.WrapUser(err, "unable to touch file")
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
			return errs.WrapUser(err, "unable to create file")
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
	dir, _ := filepath.Split(path)
	ospath := filepath.FromSlash(dir)
	err := dest.MkdirAll(ospath, 0755)
	if err != nil {
		return errs.WrapUserf(err, "couldn't create %s directory", dir)
	}

	log.Infof("%s templated", path)
	writer, err := dest.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return errs.WrapUser(err, "unable to open file")
	}
	t, e := util.OpenTemplate(sourceFile)
	if e != nil {
		return e
	}
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

func applyModuleInvocation(fs afero.Fs, path, moduleAddress string, box packr.Box) error {
	e := fs.MkdirAll(path, 0755)
	if e != nil {
		return errs.WrapUserf(e, "couldn't create %s directory", path)
	}

	moduleConfig, e := util.DownloadAndParseModule(moduleAddress)
	if e != nil {
		return errs.WrapUser(e, "could not download or parse module")
	}

	// This should really be part of the plan stage, not apply. But going to
	// leave it here for now and re-think it when we make this mechanism
	// general purpose.
	variables := make([]string, 0)
	for _, v := range moduleConfig.Variables {
		variables = append(variables, v.Name)
	}
	sort.Strings(variables)
	outputs := make([]string, 0)
	for _, o := range moduleConfig.Outputs {
		outputs = append(outputs, o.Name)
	}
	sort.Strings(outputs)
	moduleName := filepath.Base(moduleAddress)
	re := regexp.MustCompile(`\?ref=.*`)
	moduleName = re.ReplaceAllString(moduleName, "")

	moduleAddressForSource, _ := calculateModuleAddressForSource(path, moduleAddress)
	// MAIN
	f, e := box.Open("main.tf.tmpl")
	if e != nil {
		return errs.WrapUser(e, "could not open template file")
	}
	e = applyTemplate(f, fs, filepath.Join(path, "main.tf"), &moduleData{moduleName, moduleAddressForSource, variables, outputs})
	if e != nil {
		return errs.WrapUser(e, "unable to apply template for main.tf")
	}
	e = fmtHcl(fs, filepath.Join(path, "main.tf"))
	if e != nil {
		return errs.WrapUser(e, "unable to format main.tf")
	}

	// OUTPUTS
	f, e = box.Open("outputs.tf.tmpl")
	if e != nil {
		return errs.WrapUser(e, "could not open template file")
	}

	e = applyTemplate(f, fs, filepath.Join(path, "outputs.tf"), &moduleData{moduleName, moduleAddressForSource, variables, outputs})
	if e != nil {
		return errs.WrapUser(e, "unable to apply template for outputs.tf")
	}

	e = fmtHcl(fs, filepath.Join(path, "outputs.tf"))
	if e != nil {
		return errs.WrapUser(e, "unable to format outputs.tf")
	}

	return nil
}

func calculateModuleAddressForSource(path, moduleAddress string) (string, error) {
	// For cases where the module is a local path, we need to calculate the
	// relative path from the component to the module.
	// The module_source path in the fogg.json is relative to the repo root.
	var moduleAddressForSource string
	// getter will kinda normalize the module address, but it will actually be
	// wrong for local file paths, so we need to calculate that ourselves below
	s, e := getter.Detect(moduleAddress, path, getter.Detectors)
	if e != nil {
		return "", e
	}
	u, e := url.Parse(s)
	if e != nil || u.Scheme == "file" {
		// This indicates that we have a local path to the module.
		// It is possible that this test is unreliable.
		moduleAddressForSource, e = filepath.Rel(path, moduleAddress)
		if e != nil {
			return "", e
		}
	} else {
		moduleAddressForSource = moduleAddress
	}
	return moduleAddressForSource, nil
}
func getTargetPath(basePath, path string) string {
	target := filepath.Join(basePath, path)
	extension := filepath.Ext(path)

	if extension == ".tmpl" || extension == ".touch" || extension == ".create" || extension == ".rm" || extension == ".ln" {
		target = removeExtension(target)
	}

	return target
}

func linkFile(fs afero.Fs, name, target string) error {
	log.Debugf("linking %s to %s", name, target)
	linker, ok := fs.(afero.Symlinker)

	if !ok {
		return errs.NewInternal("can't cast to afero.SymLinker")
	}

	relativePath, err := filepathRel(name, target)
	log.Debugf("relative link %s err %#v", relativePath, err)
	if err != nil {
		return err
	}

	log.Debugf("removing link at %s", name)
	err = fs.Remove(name)
	log.Debugf("error removing file %s (probably ok): %s", name, err)

	_, err = linker.SymlinkIfPossible(target, name)
	return err
}

func filepathRel(path, name string) (string, error) {
	dirs := strings.Count(path, "/")
	fullPath := fmt.Sprintf("%s/%s", strings.Repeat("../", dirs), name)
	return filepath.Clean(fullPath), nil
}
