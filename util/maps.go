package util

import (
	"reflect"
	"sort"
)

func SortedMapKeys(in interface{}) []string {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Map {
		var keys []string
		for _, key := range v.MapKeys() {
			if key.Kind() == reflect.String {
				keys = append(keys, key.String())
			}
		}
		sort.Strings(keys)
		return keys
	}
	return []string{}
}
