package util

import (
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	ctyjson "github.com/zclconf/go-cty/cty/json"
	"gopkg.in/yaml.v3"
)

func dict(in interface{}) map[string]interface{} {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Map {
		r := make(map[string]interface{})
		for _, key := range v.MapKeys() {
			strct := v.MapIndex(key)
			logrus.Debug(key.Interface(), strct.Interface())
			r[key.String()] = strct.Interface()
		}
		return r
	}
	return nil
}

// https://stackoverflow.com/questions/44675087/golang-template-variable-isset
func avail(name string, data interface{}) bool {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false
	}
	return v.FieldByName(name).IsValid()
}

// https://github.com/helm/helm/blob/v3.10.1/pkg/engine/funcs.go#L79

// toYAML takes an interface, marshals it to yaml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func toYAML(v any) string {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	err := yamlEncoder.Encode(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return strings.TrimSuffix(b.String(), "\n")
}

// https://github.com/gruntwork-io/terragrunt/blob/v0.51.0/codegen/generate.go#L156

// toHCLBlock generates an HCL Block of certain type and label
// it marshals to hcl, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// {{ range $k,$v := .RequiredProviders }}
// {{ $v.Config | toHclBlock "provider" $k }}
// {{ end }}
// This is designed to be called from a template.
func toHCLBlock(blockType string, name string, config map[string]any) string {
	f := hclwrite.NewEmptyFile()
	rootBlock := f.Body().AppendNewBlock(blockType, []string{name})
	rootBlockBody := rootBlock.Body()
	var blockKeys []string

	for key := range config {
		blockKeys = append(blockKeys, key)
	}
	sort.Strings(blockKeys)
	for _, key := range blockKeys {
		// Since we don't have the cty type information for the config and since config can be arbitrary, we cheat by using
		// json as an intermediate representation.
		jsonBytes, err := json.Marshal(config[key])
		if err != nil {
			// Swallow errors inside of a template.
			return ""
		}
		var ctyVal ctyjson.SimpleJSONValue
		if err := ctyVal.UnmarshalJSON(jsonBytes); err != nil {
			// Swallow errors inside of a template.
			return ""
		}

		rootBlockBody.SetAttributeValue(key, ctyVal.Value)
	}
	return string(f.Bytes())
}

// toHCLAssignment generates an HCL assignment. It will
// always return a string, even on marshal error (empty string).
//
// {{- range $k, $v := .ProviderVersions }}
// {{ toHclAssignment $k $v }}
// {{- end }}
//
//	foo = {
//	  source  = "hashicorp/archive"
//	  version = "~> 2.0"
//	}
//
// This is designed to be called from a template.
func toHCLAssignment(name string, value any) string {
	f := hclwrite.NewEmptyFile()
	rootBlockBody := f.Body()
	// Since we don't have the cty type information for the config and since config can be arbitrary, we cheat by using
	// json as an intermediate representation.
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	var ctyVal ctyjson.SimpleJSONValue
	if err := ctyVal.UnmarshalJSON(jsonBytes); err != nil {
		// Swallow errors inside of a template.
		return ""
	}

	rootBlockBody.SetAttributeValue(name, ctyVal.Value)
	return string(f.Bytes())
}

// toHCLExpression converts a Go value to its HCL representation as a string.
// It can be used within templates to include the value in expressions or function calls.
//
// tags = merge(var.tags, {{ toHCLExpression .DefaultTags.Tags }})
//
//	tags = merge(var.tags, {
//	  Component  = "Vox"
//	  Env        = "Foo"
//	})
//
// This is designed to be called from a template.
func toHCLExpression(value any) string {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	// Since we don't have the cty type information for the value and since it can be arbitrary, we use
	// JSON as an intermediate representation.
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	var ctyVal ctyjson.SimpleJSONValue
	if err := ctyVal.UnmarshalJSON(jsonBytes); err != nil {
		// Swallow errors inside of a template.
		return ""
	}

	tokens := hclwrite.TokensForValue(ctyVal.Value)
	rootBody.AppendUnstructuredTokens(tokens)
	return strings.TrimSpace(string(f.Bytes()))
}

// deRef is a generic function to dereference a pointer to it's actual value type.
//
// This is designed to be called from a template.
func deRef[T any](v *T) T {
	return *v
}

// OpenTemplate will read `source` for a template, parse, configure and return a template.Template
func OpenTemplate(label string, source io.Reader, templates fs.FS) (*template.Template, error) {
	// TODO we should probably cache these rather than open and parse them for every apply
	funcs := sprig.TxtFuncMap()
	funcs["dict"] = dict
	funcs["avail"] = avail
	funcs["toYaml"] = toYAML
	funcs["deRefStr"] = deRef[string]
	funcs["deRefBool"] = deRef[bool]
	funcs["toHclBlock"] = toHCLBlock
	funcs["toHclAssignment"] = toHCLAssignment
	funcs["toHCLExpression"] = toHCLExpression

	s, err := io.ReadAll(source)
	if err != nil {
		return nil, errs.WrapInternal(err, "could not read template")
	}

	t, err := template.New(label).Funcs(funcs).Parse(string(s))
	if err != nil {
		return nil, errs.WrapInternal(err, "could not read template")
	}

	err = fs.WalkDir(templates, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil // nothing to do
		}

		contents, err := fs.ReadFile(templates, path)
		if err != nil {
			return errors.Wrapf(err, "could not read contents at %s", path)
		}

		t, err = t.Parse(string(contents))
		return errors.Wrapf(err, "could not parse template at %s", path)
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not walk templates")
	}
	return t, nil
}
