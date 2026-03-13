package util

import (
	"fmt"
	"io"
	"io/fs"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/pkg/errors"
)

func dict(in interface{}) map[string]interface{} {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Map {
		r := make(map[string]interface{})
		for _, key := range v.MapKeys() {
			strct := v.MapIndex(key)
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

func RenderHCLBody(config map[string]any, indent int) string {
	if len(config) == 0 {
		return ""
	}

	keys := make([]string, 0, len(config))
	for k := range config {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	prefix := strings.Repeat(" ", indent)
	for _, k := range keys {
		renderHCLEntry(&sb, k, config[k], prefix, indent)
	}
	return sb.String()
}

func renderHCLEntry(sb *strings.Builder, key string, val any, prefix string, baseIndent int) {
	switch v := val.(type) {
	case map[string]any:
		sb.WriteString(fmt.Sprintf("%s%s {\n", prefix, key))
		innerPrefix := prefix + strings.Repeat(" ", baseIndent)

		innerKeys := make([]string, 0, len(v))
		for k := range v {
			innerKeys = append(innerKeys, k)
		}
		sort.Strings(innerKeys)
		for _, ik := range innerKeys {
			renderHCLEntry(sb, ik, v[ik], innerPrefix, baseIndent)
		}
		sb.WriteString(fmt.Sprintf("%s}\n", prefix))
	case []any:
		items := make([]string, len(v))
		for i, item := range v {
			items[i] = fmt.Sprintf("%q", fmt.Sprint(item))
		}
		sb.WriteString(fmt.Sprintf("%s%s = [%s]\n", prefix, key, strings.Join(items, ", ")))
	case bool:
		sb.WriteString(fmt.Sprintf("%s%s = %t\n", prefix, key, v))
	case float64:
		if v == float64(int(v)) {
			sb.WriteString(fmt.Sprintf("%s%s = %d\n", prefix, key, int(v)))
		} else {
			sb.WriteString(fmt.Sprintf("%s%s = %g\n", prefix, key, v))
		}
	case int:
		sb.WriteString(fmt.Sprintf("%s%s = %d\n", prefix, key, v))
	default:
		sb.WriteString(fmt.Sprintf("%s%s = %q\n", prefix, key, fmt.Sprint(v)))
	}
}

// OpenTemplate will read `source` for a template, parse, configure and return a template.Template
func OpenTemplate(label string, source io.Reader, templates fs.FS) (*template.Template, error) {
	// TODO we should probably cache these rather than open and parse them for every apply
	funcs := sprig.TxtFuncMap()
	funcs["dict"] = dict
	funcs["avail"] = avail
	funcs["renderHCLBody"] = RenderHCLBody

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
