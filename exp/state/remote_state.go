package state

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/chanzuckerberg/fogg/util"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform/lang"
	"github.com/spf13/afero"
)

// liberal borrowing from https://github.com/hashicorp/terraform-config-inspect/blob/c481b8bfa41ea9dca417c2a8a98fd21bd0399e14/tfconfig/load_hcl.go#L16
func Run(_ afero.Fs, path string) error {

	fmt.Println("start")
	fs := NewOsFs()
	fmt.Printf("fs: %#v\n", fs)

	primaryPaths, diags := dirFiles(fs, path)

	parser := hclparse.NewParser()

	for _, filename := range primaryPaths {
		b, err := fs.ReadFile(filename)
		if err != nil {
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

		for _, block := range content.Blocks {

			switch block.Type {

			case "output":
				name := block.Labels[0]

				c, _, contentDiags := block.Body.PartialContent(outputSchema)
				diags = append(diags, contentDiags...)

				fmt.Println("name", name)

				if attr, defined := c.Attributes["value"]; defined {
					refs, _ := lang.ReferencesInExpr(attr.Expr)
					for _, r := range refs {
						util.Dump(r.Subject)
					}
				}

			case "module":
				name := block.Labels[0]
				fmt.Println("module", name)

				attrs, _ := block.Body.JustAttributes()

				for _, v := range attrs {
					refs, _ := lang.ReferencesInExpr(v.Expr)
					util.Dump(refs)
				}

			default:

				fmt.Println("unhandled", block.Type)
			}

		}
	}

	return nil
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
