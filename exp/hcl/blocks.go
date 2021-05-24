package hcl

import (
	"path/filepath"
	"strings"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/sirupsen/logrus"
)

func ForeachBlock(path string, f func(*hcl.Block) error) error {
	fs := tfconfig.NewOsFs()
	primaryPaths, err := dirFiles(fs, path)
	if err != nil {
		return err
	}

	parser := hclparse.NewParser()

	for _, filename := range primaryPaths {
		logrus.Debugf("reading file %s", filename)
		b, err := fs.ReadFile(filename)
		if err != nil {
			return err
		}

		var file *hcl.File
		var fileDiags hcl.Diagnostics

		if strings.HasSuffix(filename, ".json") {
			file, fileDiags = parser.ParseJSON(b, filename)
		} else {
			file, fileDiags = parser.ParseHCL(b, filename)
		}
		if fileDiags.HasErrors() {
			return fileDiags
		}

		if file == nil {
			continue
		}

		content, _, contentDiags := file.Body.PartialContent(rootSchema)
		if contentDiags.HasErrors() {
			return contentDiags
		}

		logrus.Debugf("len(content.Blocks) %v", len(content.Blocks))
		for _, block := range content.Blocks {
			err = f(block)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// taken from https://github.com/hashicorp/terraform-config-inspect/blob/c481b8bfa41ea9dca417c2a8a98fd21bd0399e14/tfconfig/load.go#L81
func dirFiles(fs tfconfig.FS, dir string) ([]string, error) {
	var primary []string

	infos, err := fs.ReadDir(dir)
	if err != nil {
		return nil, err
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

	return primary, nil
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
