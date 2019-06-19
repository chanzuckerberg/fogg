package util

import (
	"testing"

	"github.com/chanzuckerberg/fogg/templates"
	"github.com/gobuffalo/packr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDict(t *testing.T) {
	m := make(map[string]string)
	m["foo"] = "bar"
	r := dict(m)
	assert.NotNil(t, r)
	assert.IsType(t, map[string]interface{}{}, r)
	assert.Equal(t, "bar", r["foo"])
}

func TestOpenTemplate(t *testing.T) {
	temps := templates.Templates

	type args struct {
		box  packr.Box
		path string
	}
	tests := []struct {
		name    string
		args    args
		tLen    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{"foo", args{temps.Account, "Makefile.tmpl"}, 4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			f, err := tt.args.box.Open(tt.args.path)
			r.NoError(err)

			temp, err := OpenTemplate(f, &temps.Common)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			r.NotNil(temp.Templates())
			r.Len(temp.Templates(), tt.tLen)

			// if !reflect.DeepEqual(got, tt.want) {
			// t.Errorf("OpenTemplate() = %v, want %v", got, tt.want)
			// }
		})
	}
}
