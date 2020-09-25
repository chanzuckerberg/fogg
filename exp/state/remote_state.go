package state

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/chanzuckerberg/fogg/config"
	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/lang"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// liberal borrowing from https://github.com/hashicorp/terraform-config-inspect/blob/c481b8bfa41ea9dca417c2a8a98fd21bd0399e14/tfconfig/load_hcl.go#L16
func Run(fs afero.Fs, configFile, path string) error {
	// figure out which component or account we are talking about
	conf, err := config.FindAndReadConfig(fs, configFile)
	if err != nil {
		return err
	}

	componentType, componentName, envName, err := conf.PathToComponentType(path)
	if err != nil {
		return err
	}
	logrus.Debugf("componentType: %s", componentType)

	// collect remote state references
	references, err := collectRemoteStateReferences(path)
	if err != nil {
		return err
	}

	logrus.Debugf("in %s found references %#v", path, references)

	// for each reference, figure out if it is an account or component, since those are separate in our configs
	accounts := []string{}
	components := []string{}

	// we do accounts for both accounts and env components
	for _, r := range references {
		if _, found := conf.Accounts[r]; found {
			accounts = append(accounts, r)
		}
	}

	if componentType == "env" {
		env := conf.Envs[envName]

		for _, r := range references {
			if _, found := env.Components[r]; found {
				components = append(components, r)
			}
		}
	}

	// update fogg.yml with new references
	logrus.Debugf("found accounts %#v", accounts)
	logrus.Debugf("found components %#v", components)

	if componentType == "accounts" {
		if len(accounts) > 0 {
			c := conf.Accounts[componentName]
			if c.Common.DependsOn == nil {
				c.Common.DependsOn = &v2.DependsOn{}
			}

			c.DependsOn.Accounts = accounts
			conf.Accounts[componentName] = c
		}
	} else if componentType == "env" {
		if len(accounts) > 0 || len(components) > 0 {
			c := conf.Envs[envName].Components[componentName]

			if c.Common.DependsOn == nil {
				c.Common.DependsOn = &v2.DependsOn{}
			}

			c.DependsOn.Accounts = accounts
			c.DependsOn.Components = components

			conf.Envs[envName].Components[componentName] = c
		}
	}

	conf.Write(fs, configFile)

	return nil
}

func collectRemoteStateReferences(path string) ([]string, error) {
	fs := NewOsFs()

	references := map[string]bool{}

	primaryPaths, diags := dirFiles(fs, path)

	parser := hclparse.NewParser()

	for _, filename := range primaryPaths {
		logrus.Debugf("reading file %s", filename)
		b, err := fs.ReadFile(filename)
		if err != nil {
			logrus.Errorf("unable to read file %s", filename)
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Failed to read file",
				Detail:   fmt.Sprintf("The configuration file %q could not be read.", filename),
			})
			continue
		}

		var file *hcl.File
		var fileDiags hcl.Diagnostics

		if strings.HasSuffix(filename, ".json") {
			file, fileDiags = parser.ParseJSON(b, filename)
		} else {
			file, fileDiags = parser.ParseHCL(b, filename)
		}
		diags = append(diags, fileDiags...)
		if file == nil {
			continue
		}

		content, _, contentDiags := file.Body.PartialContent(rootSchema)
		diags = append(diags, contentDiags...)

		logrus.Debugf("len(content.Blocks) %v", len(content.Blocks))
		for _, block := range content.Blocks {
			logrus.Debugf("block type: %v", block.Type)

			attrs, _ := block.Body.JustAttributes()

			for _, v := range attrs {
				refs, _ := lang.ReferencesInExpr(v.Expr)

				for _, r := range refs {
					if r != nil {
						logrus.Debugf("ref: %v", r)
						if resource, ok := r.Subject.(addrs.ResourceInstance); ok {
							if resource.Resource.Type == "terraform_remote_state" {
								references[resource.Resource.Name] = true
							}
						}
					}
				}
			}
		}
	}

	// FIXME return diags and make it work with error reporting
	var refNames []string
	for k := range references {
		refNames = append(refNames, k)
	}
	return refNames, nil
}

// taken from https://github.com/hashicorp/terraform-config-inspect/blob/c481b8bfa41ea9dca417c2a8a98fd21bd0399e14/tfconfig/load.go#L81
func dirFiles(fs FS, dir string) (primary []string, diags hcl.Diagnostics) {
	infos, err := fs.ReadDir(dir)
	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to read module directory",
			Detail:   fmt.Sprintf("Module directory %s does not exist or cannot be read.", dir),
		})
		return
	}

	var override []string
	for _, info := range infos {
		if info.IsDir() {
			// We only care about files
			continue
		}

		name := info.Name()
		ext := fileExt(name)
		if ext == "" || isIgnoredFile(name) {
			continue
		}

		baseName := name[:len(name)-len(ext)] // strip extension
		isOverride := baseName == "override" || strings.HasSuffix(baseName, "_override")

		fullPath := filepath.Join(dir, name)
		if isOverride {
			override = append(override, fullPath)
		} else {
			primary = append(primary, fullPath)
		}
	}

	// We are assuming that any _override files will be logically named,
	// and processing the files in alphabetical order. Primaries first, then overrides.
	primary = append(primary, override...)

	return
}

// fileExt returns the Terraform configuration extension of the given
// path, or a blank string if it is not a recognized extension.
func fileExt(path string) string {
	if strings.HasSuffix(path, ".tf") {
		return ".tf"
	} else if strings.HasSuffix(path, ".tf.json") {
		return ".tf.json"
	} else {
		return ""
	}
}

// isIgnoredFile returns true if the given filename (which must not have a
// directory path ahead of it) should be ignored as e.g. an editor swap file.
func isIgnoredFile(name string) bool {
	return strings.HasPrefix(name, ".") || // Unix-like hidden files
		strings.HasSuffix(name, "~") || // vim
		strings.HasPrefix(name, "#") && strings.HasSuffix(name, "#") // emacs
}

// taken from https://github.com/hashicorp/terraform-config-inspect/blob/c481b8bfa41ea9dca417c2a8a98fd21bd0399e14/tfconfig/filesystem.go

// FS represents a minimal filesystem implementation
// See io/fs.FS in http://golang.org/s/draft-iofs-design
type FS interface {
	Open(name string) (File, error)
	ReadFile(name string) ([]byte, error)
	ReadDir(dirname string) ([]os.FileInfo, error)
}

// File represents an open file in FS
// See io/fs.File in http://golang.org/s/draft-iofs-design
type File interface {
	Stat() (os.FileInfo, error)
	Read([]byte) (int, error)
	Close() error
}

type osFs struct{}

func (fs *osFs) Open(name string) (File, error) {
	return os.Open(name)
}

func (fs *osFs) ReadFile(name string) ([]byte, error) {
	return ioutil.ReadFile(name)
}

func (fs *osFs) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

// NewOsFs provides a basic implementation of FS for an OS filesystem
func NewOsFs() FS {
	return &osFs{}
}
