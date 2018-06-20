package util

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/gobuffalo/packr"
)

func dict(in interface{}) map[string]interface{} {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Map {
		r := make(map[string]interface{})
		for _, key := range v.MapKeys() {
			strct := v.MapIndex(key)
			fmt.Println(key.Interface(), strct.Interface())
			r[key.String()] = strct.Interface()
		}
		return r
	}
	return nil
}

func OpenTemplate(source packr.File) *template.Template {
	s, err := ioutil.ReadAll(source)
	if err != nil {
		panic(err)
	}
	funcs := sprig.TxtFuncMap()
	funcs["dict"] = dict
	return template.Must(template.New("tmpl").Funcs(funcs).Parse(string(s)))
}
