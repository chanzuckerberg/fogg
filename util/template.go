package util

import (
	"io"
	"io/ioutil"
	"reflect"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/chanzuckerberg/fogg/errs"
	"github.com/gobuffalo/packr"
	"github.com/sirupsen/logrus"
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

// https://gist.github.com/morrxy/06fb44e967ba5c8b4b59f08268f9dd73
func hasField(v interface{}, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByName(name).IsValid()
}

// OpenTemplate will read `source` for a template, parse, configure and return a template.Template
func OpenTemplate(path string, source io.Reader, commonTemplates *packr.Box) (*template.Template, error) {
	// TODO we should probably cache these rather than open and parse them for every apply

	var readTemplate = func(source io.Reader) (string, error) {
		s, err := ioutil.ReadAll(source)

		if err != nil {
			return "", errs.WrapInternal(err, "could not read template")
		}
		return string(s), nil
	}

	s, err := readTemplate(source)
	if err != nil {
		return nil, err
	}

	funcs := sprig.TxtFuncMap()
	funcs["dict"] = dict
	funcs["hasField"] = hasField

	t, err := template.New(path).Funcs(funcs).Parse(s)
	if err != nil {
		return nil, err
	}

	err = commonTemplates.Walk(func(path string, file packr.File) error {
		logrus.Debugf("parsing common template %s", path)
		s, err := readTemplate(file)
		if err != nil {
			return err
		}

		t, err = t.Parse(s)

		return err

	})

	if err != nil {
		return nil, err
	}

	return t, err
}
