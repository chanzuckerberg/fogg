package config

import (
	"testing"

	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func getPtr(val string) *string {
	return &val
}

func TestInitConfig(t *testing.T) {
	r := require.New(t)
	c := InitConfig(getPtr("proj"), getPtr("reg"), getPtr("buck"), getPtr("table"), getPtr("prof"), getPtr("me@foo.example"), getPtr("123456789"), "0.99.0")
	r.Equal("prof", *c.Defaults.Common.Backend.Profile)
	r.Equal("prof", *c.Defaults.Providers.AWS.Profile)
	r.Equal("reg", *c.Defaults.Providers.AWS.Region)
	r.Equal("reg", *c.Defaults.Providers.AWS.Region)
	r.Equal("0.99.0", *c.Defaults.Providers.AWS.Version)
	r.Equal("buck", *c.Defaults.Common.Backend.Bucket)
	r.Equal("table", *c.Defaults.Common.Backend.DynamoTable)
	r.Equal("me@foo.example", *c.Defaults.Owner)
	r.Equal("proj", *c.Defaults.Project)
	r.Equal(defaultTerraformVersion.String(), *c.Defaults.TerraformVersion)
}

func Test_detectVersion(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"explicit 2", args{[]byte(`{"version": 2}`)}, 2, false},
		{"err", args{[]byte(`{`)}, 0, true},
	}
	r := require.New(t)
	fs, _, err := util.TestFs()
	r.NoError(err)

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			var got int
			var err error
			switch tt.want {
			case 2:
				afero.WriteFile(fs, "fogg.yml", tt.args.b, 0644) //nolint:errcheck
				got, err = detectVersion(tt.args.b, fs, "fogg.yml")
			default:
				got, err = detectVersion(tt.args.b, fs, "fogg.yml")
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("detectVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("detectVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
