package util

import (
	"io"
	"io/ioutil"
	"reflect"
	"text/template"

	"github.com/Masterminds/sprig"
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

func OpenTemplate(source io.Reader) *template.Template {
	s, err := ioutil.ReadAll(source)
	if err != nil {
		log.Panic(err) // FIXME
	}
	funcs := sprig.TxtFuncMap()
	funcs["dict"] = dict
	return template.Must(template.New("tmpl").Funcs(funcs).Parse(string(s)))
}
