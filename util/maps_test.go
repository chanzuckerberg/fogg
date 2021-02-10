package util

import (
	"reflect"
	"testing"
)

func TestSortedMapKeys(t *testing.T) {
	type args struct {
		in interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"basic", args{map[string]string{"foo": "bar"}}, []string{"foo"}},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			if got := SortedMapKeys(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortedMapKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
