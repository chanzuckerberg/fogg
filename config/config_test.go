package config

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	a := assert.New(t)
	c := InitConfig("proj", "reg", "buck", "table", "prof", "me@foo.example", "0.99.0")
	a.Equal("prof", *c.Defaults.Common.Backend.Profile)
	a.Equal("prof", *c.Defaults.Providers.AWS.Profile)
	a.Equal("reg", *c.Defaults.Providers.AWS.Region)
	a.Equal("reg", *c.Defaults.Providers.AWS.Region)
	a.Equal("0.99.0", *c.Defaults.Providers.AWS.Version)
	a.Equal("buck", *c.Defaults.Common.Backend.Bucket)
	a.Equal("table", *c.Defaults.Common.Backend.DynamoTable)
	a.Equal("me@foo.example", *c.Defaults.Owner)
	a.Equal("proj", *c.Defaults.Project)
	a.Equal(defaultTerraformVersion.String(), *c.Defaults.TerraformVersion)
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
		{"implicit 1", args{[]byte(`{}`)}, 1, false},
		{"explicit 1", args{[]byte(`{"version": 1}`)}, 1, false},
		{"explicit 2", args{[]byte(`{"version": 2}`)}, 2, false},
		{"err", args{[]byte(`{`)}, 0, true},
	}
	a := assert.New(t)
	fs, _, err := util.TestFs()
	a.NoError(err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got int
			var err error
			switch tt.want {
			case 1:
				afero.WriteFile(fs, "fogg.json", tt.args.b, 0644)
				got, err = detectVersion(tt.args.b, fs, "fogg.json")
			case 2:
				afero.WriteFile(fs, "fogg.yml", tt.args.b, 0644)
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

func intptr(i int) *int {
	return &i
}

func jsonNumberPtr(i int) *json.Number {
	j := json.Number(strconv.Itoa(i))
	return &j
}

func strptr(s string) *string {
	return &s
}

func boolptr(b bool) *bool {
	return &b
}

func TestFindAndReadConfig(t *testing.T) {
	a := assert.New(t)

	fs := func(m map[string][]byte) (afero.Fs, error) {
		fs, _, err := util.TestFs()
		a.NoError(err)

		for k, v := range m {
			f, e := fs.OpenFile(k, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if e != nil {
				return nil, e
			}
			_, e = f.Write(v)
			if e != nil {
				return nil, e
			}
		}
		return fs, nil
	}

	v2, e := util.TestFile("v2_minimal_valid")

	a.NoError(e)

	f2, e := fs(map[string][]byte{
		"config.json": v2,
	})
	a.NoError(e)
	defer f2.RemoveAll(".") //nolint

	fErr, e := fs(map[string][]byte{
		"config.json": []byte(`{"version": 7}`),
	})
	a.NoError(e)
	defer fErr.RemoveAll(".") //nolint

	_, e = FindAndReadConfig(f2, "config.json")
	a.NoError(e)

	_, e = FindAndReadConfig(fErr, "config.json")
	a.Error(e)

}
