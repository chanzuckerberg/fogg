package apply

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/plan"
	"github.com/chanzuckerberg/fogg/templates"
	"github.com/chanzuckerberg/fogg/util"
	getter "github.com/hashicorp/go-getter"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const rootPath = "terraform"

// Apply will run a plan and apply all the changes to the current repo.
func Apply(fs afero.Fs, conf *v2.Config, tmp *templates.T, upgrade bool) error {
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
	plan, err := plan.Eval(conf)
	if err != nil {
		return errs.WrapUser(err, "unable to evaluate plan")
	}
	err = applyRepo(fs, plan, tmp.Repo, tmp.Common)
	if err != nil {
		return errs.WrapUser(err, "unable to apply repo")
	}

	if plan.TravisCI.Enabled {
		err = applyTree(fs, tmp.TravisCI, tmp.Common, "", plan.TravisCI)
		if err != nil {
			return errs.WrapUser(err, "unable to apply travis ci")
		}
	}

	if plan.CircleCI.Enabled {
		err = applyTree(fs, tmp.CircleCI, tmp.Common, "", plan.CircleCI)
		if err != nil {
			return errs.WrapUser(err, "unable to apply CircleCI")
		}
	}

	if plan.GitHubActionsCI.Enabled {
		err = applyTree(fs, tmp.GitHubActionsCI, tmp.Common, ".github", plan.GitHubActionsCI)
		if err != nil {
			return errs.WrapUser(err, "unable to apply GitHub Actions CI")
		}
	}

	tfBox := tmp.Components[v2.ComponentKindTerraform]
	err = applyAccounts(fs, plan, tfBox, tmp.Common)
	if err != nil {
		return errs.WrapUser(err, "unable to apply accounts")
	}

	err = applyEnvs(fs, plan, tmp.Env, tmp.Components, tmp.Common)
	if err != nil {
		return errs.WrapUser(err, "unable to apply envs")
	}

	tfBox = tmp.Components[v2.ComponentKindTerraform]
	err = applyGlobal(fs, plan.Global, tfBox, tmp.Common)
	if err != nil {
		return errs.WrapUser(err, "unable to apply global")
	}

	err = applyModules(fs, plan.Modules, tmp.Module, tmp.Common)
	if err != nil {
		return errs.WrapUser(err, "unable to apply modules")
	}

	return applyTFE(fs, plan)
}

type LocalsTFE struct {
	Locals *Locals `json:"locals,omitempty"`
}

type Locals struct {
	Accounts         map[string]*TFEWorkspace            `json:"accounts,omitempty"`
	Envs             map[string]map[string]*TFEWorkspace `json:"envs,omitempty"`
	DefaultTFVersion *string                             `json:"default_terraform_version,omitempty"`
}
type ExtraTeamPermissions struct {
	Plan *[]string `json:"plan,omitempty"`
	Read *[]string `json:"read,omitempty"`
}
type TFEWorkspace struct {
	TriggerPrefixes      *[]string             `json:"trigger_prefixes"`
	WorkingDirectory     *string               `json:"working_directory,omitempty"`
	TerraformVersion     *string               `json:"terraform_version,omitempty"`
	ExtraTeamPermissions *ExtraTeamPermissions `json:"extra_team_permissions,omitempty"`
	GithubBranch         *string               `json:"branch,omitempty"`
	AutoApply            *bool                 `json:"auto_apply,omitempty"`
	RemoteApply          *bool                 `json:"remote_apply,omitempty"`
}

func MakeTFEWorkspace() *TFEWorkspace {
	defaultTriggerPrefixes := make([]string, 0)
	defaultTerraformVersion := "1.2.6"
	defaultGithubBranch := "main"
	defaultAutoApply := true
	defaultRemoteApply := false
	return &TFEWorkspace{
		TriggerPrefixes:  &defaultTriggerPrefixes,
		TerraformVersion: &defaultTerraformVersion,
		GithubBranch:     &defaultGithubBranch,
		AutoApply:        &defaultAutoApply,
		RemoteApply:      &defaultRemoteApply,
	}
}

func updateLocalsFromPlan(locals *LocalsTFE, plan *plan.Plan) {
	// if there is a planned env or account that isn't in the locals, add it
	for accountName := range plan.Accounts {
		if _, ok := locals.Locals.Accounts[accountName]; !ok {
			locals.Locals.Accounts[accountName] = MakeTFEWorkspace()
		}
	}
	for envName := range plan.Envs {
		if _, ok := locals.Locals.Envs[envName]; !ok {
			locals.Locals.Envs[envName] = make(map[string]*TFEWorkspace, 0)
		}
		for componentName := range plan.Envs[envName].Components {
			if _, ok := locals.Locals.Envs[envName][componentName]; !ok {
				locals.Locals.Envs[envName][componentName] = MakeTFEWorkspace()
			}
		}
	}

	// if there is a locals env or account that isn't in the plan, delete it
	for account := range locals.Locals.Accounts {
		shouldDelete := func() bool {
			for plannedAccount := range plan.Accounts {
				if account == plannedAccount {
					return false
				}
			}
			return true
		}()
		if shouldDelete {
			delete(locals.Locals.Accounts, account)
		}
	}
	for envName, component := range locals.Locals.Envs {
		for componentName := range component {
			shouldDelete := func() bool {
				for plannedComponent := range plan.Envs[envName].Components {
					if plannedComponent == componentName {
						return false
					}
				}

				return true
			}()
			if shouldDelete {
				delete(locals.Locals.Envs[envName], componentName)
			}
		}
	}
}

func applyTFE(fs afero.Fs, plan *plan.Plan) error {
	tfePath := filepath.Join("terraform", "tfe", "locals.tf.json")
	_, err := fs.Stat(tfePath)
	// if the repo doesn't have a locals.tf.json, don't worry about it
	if errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}
	read, err := fs.Open(tfePath)
	if err != nil {
		return errors.Wrapf(err, "unable to open locals.tf.json file %s for unmarshalling", tfePath)
	}
	defer read.Close()
	locals := LocalsTFE{}
	err = json.NewDecoder(read).Decode(&locals)
	if err != nil {
		return errors.Wrapf(err, "unable to decode locals.tf.json from %s", tfePath)
	}

	updateLocalsFromPlan(&locals, plan)
	write, err := fs.OpenFile(tfePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return errors.Wrapf(err, "unable to open locals.tf.json file %s for marshaling", tfePath)
	}
	defer write.Close()
	encoder := json.NewEncoder(write)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(locals)
	if err != nil {
		return errors.Wrap(err, "unable to marhsal locals.tf.json")
	}

	return nil
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
	repoVersion := strings.TrimSpace(string(b))
	changed := versionIsChanged(repoVersion, current)
	return changed, repoVersion, nil
}

func versionIsChanged(repo string, tool string) bool {
	repoVersion, repoSha, repoDirty := util.ParseVersion(repo)
	toolVersion, toolSha, toolDirty := util.ParseVersion(tool)

	if repoDirty || toolDirty {
		return true
	}

	if (repoSha != "" || toolSha != "") && repoSha != toolSha {
		return true
	}

	return toolVersion.NE(repoVersion)
}

func applyRepo(fs afero.Fs, p *plan.Plan, repoTemplates, commonTemplates fs.FS) error {
	return applyTree(fs, repoTemplates, commonTemplates, "", p)
}

func applyGlobal(fs afero.Fs, p plan.Component, repoBox, commonBox fs.FS) error {
	logrus.Debug("applying global")
	path := fmt.Sprintf("%s/global", rootPath)
	e := fs.MkdirAll(path, 0755)
	if e != nil {
		return errs.WrapUserf(e, "unable to make directory %s", path)
	}
	return applyTree(fs, repoBox, commonBox, path, p)
}

func applyAccounts(fs afero.Fs, p *plan.Plan, accountBox, commonBox fs.FS) (e error) {
	for account, accountPlan := range p.Accounts {
		path := fmt.Sprintf("%s/accounts/%s", rootPath, account)
		e = fs.MkdirAll(path, 0755)
		if e != nil {
			return errs.WrapUser(e, "unable to make directories for accounts")
		}
		e = applyTree(fs, accountBox, commonBox, path, accountPlan)
		if e != nil {
			return errs.WrapUser(e, "unable to apply templates to account")
		}
	}
	return nil
}

func applyModules(fs afero.Fs, p map[string]plan.Module, moduleBox, commonBox fs.FS) (e error) {
	for module, modulePlan := range p {
		path := fmt.Sprintf("%s/modules/%s", rootPath, module)
		e = fs.MkdirAll(path, 0755)
		if e != nil {
			return errs.WrapUserf(e, "unable to make path %s", path)
		}
		e = applyTree(fs, moduleBox, commonBox, path, modulePlan)
		if e != nil {
			return errs.WrapUser(e, "unable to apply tree")
		}
	}
	return nil
}

func applyEnvs(
	fs afero.Fs,
	p *plan.Plan,
	envBox fs.FS,
	componentBoxes map[v2.ComponentKind]fs.FS,
	commonBox fs.FS) (err error) {
	logrus.Debug("applying envs")
	for env, envPlan := range p.Envs {
		logrus.Debugf("applying %s", env)
		path := fmt.Sprintf("%s/envs/%s", rootPath, env)
		err = fs.MkdirAll(path, 0755)
		if err != nil {
			return errs.WrapUserf(err, "unable to make directory %s", path)
		}
		err := applyTree(fs, envBox, commonBox, path, envPlan)
		if err != nil {
			return errs.WrapUser(err, "unable to apply templates to env")
		}
		for component, componentPlan := range envPlan.Components {
			path = fmt.Sprintf("%s/envs/%s/%s", rootPath, env, component)
			err = fs.MkdirAll(path, 0755)
			if err != nil {
				return errs.WrapUser(err, "unable to make directories for component")
			}

			// NOTE(el): component kind only support TF now
			// 					 add a dynamic check to make sure.
			kind := componentPlan.Kind.GetOrDefault()
			componentBox, ok := componentBoxes[kind]
			if !ok {
				return errs.NewUserf("component of kind '%s' not suppoerted, must be 'terraform'", kind)
			}

			err := applyTree(fs, componentBox, commonBox, path, componentPlan)
			if err != nil {
				return errs.WrapUser(err, "unable to apply templates for component")
			}

			if componentPlan.ModuleSource != nil {
				downloader, err := util.MakeDownloader(*componentPlan.ModuleSource)
				if err != nil {
					return errs.WrapUser(err, "unable to make a downloader")
				}
				err = applyModuleInvocation(fs, path, *componentPlan.ModuleSource, componentPlan.ModuleName, templates.Templates.ModuleInvocation, commonBox, downloader)
				if err != nil {
					return errs.WrapUser(err, "unable to apply module invocation")
				}
			}
		}
	}
	return nil
}

func applyTree(dest afero.Fs, source fs.FS, common fs.FS, targetBasePath string, subst interface{}) (e error) {
	return fs.WalkDir(source, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errs.WrapInternal(err, "unable to walk dir")
		}
		if d.IsDir() {
			return nil // skip dirs
		}

		sourceFile, err := source.Open(path)
		if err != nil {
			return errs.WrapInternal(err, "could not read source file")
		}

		extension := filepath.Ext(path)
		target := getTargetPath(targetBasePath, path)
		targetExtension := filepath.Ext(target)

		if extension == ".tmpl" {
			e = applyTemplate(sourceFile, common, dest, target, subst)
			if e != nil {
				return errs.WrapUser(e, "unable to apply template")
			}

			if targetExtension == ".tf" {
				e = fmtHcl(dest, target, true)
				if e != nil {
					return errs.WrapUser(e, "unable to format HCL")
				}
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
			logrus.Infof("%s removed", target)
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
			logrus.Infof("%s copied", target)
		}

		return nil
	})
}

// collapseLines will convert \n+ to \n to reduce spurious diffs in generated output
// post 0.12 terraform fmt will not manage vertical whitespace
// https://github.com/hashicorp/terraform/issues/23223#issuecomment-547519852
func collapseLines(in []byte) []byte {
	fmtRegex := regexp.MustCompile(`\n+`)
	return fmtRegex.ReplaceAll(in, []byte("\n"))
}

func fmtHcl(fs afero.Fs, path string, collapse bool) error {
	in, e := afero.ReadFile(fs, path)
	if e != nil {
		return errs.WrapUserf(e, "unable to read file %s", path)
	}
	if collapse {
		in = collapseLines(in)
	}
	out := hclwrite.Format(in)
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
		logrus.Infof("%s touched", path)
		_, err = dest.Create(path)
		if err != nil {
			return errs.WrapUser(err, "unable to touch file")
		}
	} else {
		logrus.Debugf("%s skipped touch", path)
	}
	return nil
}

func createFile(dest afero.Fs, path string, sourceFile io.Reader) error {
	_, err := dest.Stat(path)
	if err != nil { // TODO we might not want to do this for all errors
		logrus.Infof("%s created", path)
		err = afero.WriteReader(dest, path, sourceFile)
		if err != nil {
			return errs.WrapUser(err, "unable to create file")
		}
	} else {
		logrus.Infof("%s skipped", path)
	}
	return nil
}

func removeExtension(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}

func applyTemplate(sourceFile io.Reader, commonTemplates fs.FS, dest afero.Fs, path string, overrides interface{}) error {
	dir, _ := filepath.Split(path)
	ospath := filepath.FromSlash(dir)
	err := dest.MkdirAll(ospath, 0775)
	if err != nil {
		return errs.WrapUserf(err, "couldn't create %s directory", dir)
	}

	logrus.Infof("%s templated", path)
	writer, err := dest.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errs.WrapUser(err, "unable to open file")
	}
	t, e := util.OpenTemplate(path, sourceFile, commonTemplates)
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

func applyModuleInvocation(
	fs afero.Fs,
	path, moduleAddress string,
	inModuleName *string,
	box fs.FS,
	commonBox fs.FS,
	downloadFunc util.ModuleDownloader,
) error {
	e := fs.MkdirAll(path, 0755)
	if e != nil {
		return errs.WrapUserf(e, "couldn't create %s directory", path)
	}

	moduleConfig, e := downloadFunc.DownloadAndParseModule(fs)
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

	moduleName := ""
	if inModuleName != nil {
		moduleName = *inModuleName
	}
	if moduleName == "" {
		moduleName = filepath.Base(moduleAddress)
		re := regexp.MustCompile(`\?ref=.*`)
		moduleName = re.ReplaceAllString(moduleName, "")
	}

	moduleAddressForSource, _ := calculateModuleAddressForSource(path, moduleAddress)
	// MAIN
	f, e := box.Open("main.tf.tmpl")
	if e != nil {
		return errs.WrapUser(e, "could not open template file")
	}
	e = applyTemplate(
		f,
		commonBox,
		fs,
		filepath.Join(path, "main.tf"),
		&moduleData{moduleName, moduleAddressForSource, variables, outputs})
	if e != nil {
		return errs.WrapUser(e, "unable to apply template for main.tf")
	}
	e = fmtHcl(fs, filepath.Join(path, "main.tf"), false)
	if e != nil {
		return errs.WrapUser(e, "unable to format main.tf")
	}

	// OUTPUTS
	f, e = box.Open("outputs.tf.tmpl")
	if e != nil {
		return errs.WrapUser(e, "could not open template file")
	}

	e = applyTemplate(f, commonBox, fs, filepath.Join(path, "outputs.tf"), &moduleData{moduleName, moduleAddressForSource, variables, outputs})
	if e != nil {
		return errs.WrapUser(e, "unable to apply template for outputs.tf")
	}

	e = fmtHcl(fs, filepath.Join(path, "outputs.tf"), false)
	if e != nil {
		return errs.WrapUser(e, "unable to format outputs.tf")
	}

	return nil
}

func calculateModuleAddressForSource(path, moduleAddress string) (string, error) {
	// For cases where the module is a local path, we need to calculate the
	// relative path from the component to the module.
	// The module_source path in the fogg.yml is relative to the repo root.
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
	logrus.Debugf("linking %s to %s", name, target)
	linker, ok := fs.(afero.Symlinker)

	if !ok {
		return errs.NewInternal("can't cast to afero.SymLinker")
	}

	relativePath := filepathRel(name, target)
	logrus.Debugf("relative link %s", relativePath)

	logrus.Debugf("removing link at %s", name)
	err := fs.Remove(name)
	logrus.Debugf("error removing file %s (probably ok): %s", name, err)

	_, err = linker.SymlinkIfPossible(relativePath, name)
	return err
}

func filepathRel(path, name string) string {
	dirs := strings.Count(path, "/")
	fullPath := fmt.Sprintf("%s/%s", strings.Repeat("../", dirs), name)
	return filepath.Clean(fullPath)
}
