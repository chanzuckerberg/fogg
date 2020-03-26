package templates

import (
	"testing"

	"github.com/chanzuckerberg/fogg/util"
	"github.com/gobuffalo/packr"
	"github.com/stretchr/testify/require"
)

func TestOpenTemplate(t *testing.T) {
	temps := Templates

	type args struct {
		box  packr.Box
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"foo", args{temps.Account, "Makefile.tmpl"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			f, err := tt.args.box.Open(tt.args.path)
			r.NoError(err)

			temp, err := util.OpenTemplate("foo", f, &temps.Common)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			r.NotNil(temp.Templates())
		})
	}
}
