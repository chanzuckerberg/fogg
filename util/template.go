package util

import (
	"io"
	"io/ioutil"
	"reflect"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/chanzuckerberg/fogg/errs"
	log "github.com/sirupsen/logrus"
)

func dict(in interface{}) map[string]interface{} {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Map {
		r := make(map[string]interface{})
		for _, key := range v.MapKeys() {
			strct := v.MapIndex(key)
			log.Debug(key.Interface(), strct.Interface())
			r[key.String()] = strct.Interface()
		}
		return r
	}
	return nil
}

// OpenTemplate will read `source` for a template, parse, configure and return a template.Template
func OpenTemplate(source io.Reader) (*template.Template, error) {
	s, err := ioutil.ReadAll(source)
	if err != nil {
		return nil, errs.WrapInternal(err, "could not read template")
	}
	funcs := sprig.TxtFuncMap()
	funcs["dict"] = dict
	return template.New("tmpl").Funcs(funcs).Parse(string(s))
}
