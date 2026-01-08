package util

import (
	"io"
	"io/fs"
	"reflect"
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

// OpenTemplate will read `source` for a template, parse, configure and return a template.Template
func OpenTemplate(label string, source io.Reader, templates fs.FS) (*template.Template, error) {
	// TODO we should probably cache these rather than open and parse them for every apply
	funcs := sprig.TxtFuncMap()
	funcs["dict"] = dict
	funcs["avail"] = avail

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
