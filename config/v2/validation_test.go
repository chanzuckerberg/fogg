package v2

import (
	"testing"

	"github.com/chanzuckerberg/fogg/util"
	"github.com/jinzhu/copier"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_nonEmptyString(t *testing.T) {
	a := assert.New(t)

	empty := ""
	nonEmpty := "foo"
	a.True(nonEmptyString(&nonEmpty))
	a.False(nonEmptyString(&empty))
	a.False(nonEmptyString(nil))
}

func TestValidateOwnersAccount(t *testing.T) {
	// this will serve as a test for all the fuctions that use validateInheritedStringField, since they are equivalent

	a := assert.New(t)
	foo := "foo@example.com"

	// acct owner

	c := confAcctOwner(foo, foo)

	// Both defaults and acct are set
	a.Nil(c.validateInheritedStringField("owner", OwnerGetter, nonEmptyString).ErrorOrNil())

	// defaults unset, still valid
	c = confAcctOwner("", foo)
	a.NoError(c.validateInheritedStringField("owner", OwnerGetter, nonEmptyString).ErrorOrNil())

	// both unset, no longer valid
	c = confAcctOwner("", "")
	a.Equal(2, c.validateInheritedStringField("owner", OwnerGetter, nonEmptyString).Len())
}

func TestValidateOwnersComponent(t *testing.T) {
	foo := "foo@example.com"

	var cases = []struct {
		label  string
		def    string
		env    string
		comp   string
		errNil bool
		errz   int
	}{
		{"all set", foo, foo, foo, true, 0},
		{"def unset", "", foo, foo, true, 0},
		{"env unset", foo, "", foo, true, 0},
		{"comp unset", foo, foo, "", true, 0},
		{"def & env unset", "", "", foo, true, 0},
		{"def & comp unset", "", foo, "", true, 0},
		{"all unset", "", "", "", false, 1},
	}

	for _, tt := range cases {
		t.Run(tt.label, func(t *testing.T) {
			a := assert.New(t)
			c := confComponentOwner(tt.def, tt.env, tt.comp)
			e := c.validateInheritedStringField("owner", OwnerGetter, nonEmptyString)
			if tt.errNil {
				a.NoError(e.ErrorOrNil())
			} else {
				a.NotNil(e)
				a.Equal(tt.errz, e.Len())
			}
		})
	}
}

func TestValidateBackends(t *testing.T) {

	var cases = []struct {
		kind     string
		genValid bool
		wantErr  bool
	}{
		// {"invalid", false, true},
		// {"s3", true, false},
		{"remote", true, false},
	}

	for _, tt := range cases {
		t.Run(tt.kind, func(t *testing.T) {
			r := require.New(t)

			c := confBackendKind(t, tt.kind, tt.genValid)
			_, err := c.Validate()
			if tt.wantErr {
				r.Error(err)
			} else {
				r.NoError(err)
			}
		})
	}
}

func confBackendKind(t *testing.T, kind string, generateValid bool) Config {
	r := require.New(t)
	base := Config{
		Version: 2,
		Defaults: Defaults{
			Common{
				Owner: util.StrPtr("foo@example.com"),
				Backend: &Backend{
					Kind:    util.StrPtr(kind),
					Bucket:  util.StrPtr("foo"),
					Region:  util.StrPtr("foo"),
					Profile: util.StrPtr("foo"),
				},
				Project:          util.StrPtr("foo"),
				TerraformVersion: util.StrPtr("1.1.1"),
			},
		},
	}

	if kind == "remote" {
		var remote Config
		err := copier.Copy(&remote, &base)
		r.NoError(err)

		remote.Defaults.Common.Backend = &Backend{
			Kind:         util.StrPtr("remote"),
			HostName:     util.StrPtr("example.com"),
			Organization: util.StrPtr("example"),
		}
		return remote
	}

	return base
}

func confAcctOwner(def, acct string) Config {
	return Config{
		Defaults: Defaults{
			Common{
				Owner: &def,
				Backend: &Backend{
					Bucket:  util.StrPtr("foo"),
					Region:  util.StrPtr("foo"),
					Profile: util.StrPtr("foo"),
				},
			},
		},
		Accounts: map[string]Account{
			"foo": Account{
				Common{
					Owner: util.StrPtr(acct),
				},
			},
		},
		Global: Component{
			Common: Common{
				Owner: util.StrPtr(acct),
			},
		},
	}
}

func confComponentOwner(def, env, component string) Config {
	return Config{
		Defaults: Defaults{
			Common{
				Owner: util.StrPtr(def),
				Backend: &Backend{
					Bucket:  util.StrPtr("foo"),
					Region:  util.StrPtr("foo"),
					Profile: util.StrPtr("foo"),
				},
			},
		},
		Envs: map[string]Env{
			"bar": Env{
				Common: Common{
					Owner: util.StrPtr(env),
				},
				Components: map[string]Component{
					"bam": Component{
						Common: Common{
							Owner: util.StrPtr(component),
						},
					},
				},
			},
		},
		Global: Component{
			Common: Common{
				Owner: util.StrPtr("foo"),
			},
		},
	}
}

func TestResolveStringArray(t *testing.T) {
	def := []string{"foo"}
	override := []string{"bar"}

	result := ResolveStringArray(def, override)
	assert.Len(t, result, 1)
	assert.Equal(t, "bar", result[0])

	override = nil

	result2 := ResolveStringArray(def, override)
	assert.Len(t, result2, 1)
	assert.Equal(t, "foo", result2[0])
}

func TestConfig_ValidateAWSProviders(t *testing.T) {

	tests := []struct {
		fileName string
		wantErr  bool
	}{
		{"v2_full_yaml", false},
		{"v2_minimal_valid_yaml", false},
		{"v2_invalid_aws_provider_yaml", true},
	}
	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			r := require.New(t)
			fs, _, e := util.TestFs()
			r.NoError(e)

			b, e := util.TestFile(tt.fileName)
			r.NoError(e)
			e = afero.WriteFile(fs, "fogg.yml", b, 0644)
			r.NoError(e)

			c, e := ReadConfig(fs, b, "fogg.yml")
			r.NoError(e)
			r.NotNil(c)

			if err := c.ValidateAWSProviders(); (err != nil) != tt.wantErr {
				t.Errorf("Config.ValidateAWSProviders() error = %v, wantErr %v (err != nil) %v", err, tt.wantErr, (err != nil))
			}
		})
	}
}

func TestConfig_ValidateTravis(t *testing.T) {
	tests := []struct {
		fileName string
		wantErr  bool
	}{
		{"v2_invalid_travis_command_yaml", true},
	}
	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			r := require.New(t)
			fs, _, e := util.TestFs()
			r.NoError(e)

			b, e := util.TestFile(tt.fileName)
			r.NoError(e)
			e = afero.WriteFile(fs, "fogg.json", b, 0644)
			r.NoError(e)

			c, e := ReadConfig(fs, b, "fogg.json")
			r.NoError(e)

			if err := c.ValidateTravis(); (err != nil) != tt.wantErr {
				t.Errorf("Config.ValidateTravis() error = %v, wantErr %v (err != nil) %v", err, tt.wantErr, (err != nil))
			}
		})
	}
}

func TestFailOnTF11(t *testing.T) {
	r := require.New(t)

	conf := confAcctOwner("foo", "bar")
	conf.Modules = map[string]Module{}

	conf.Modules["bad_tf_version"] = Module{
		TerraformVersion: util.StrPtr("0.11.0"),
	}

	err := conf.validateModules()
	r.Error(err, "fogg only supports tf versions >= 0.12.0 but 0.11.0 was provided")
}

func TestValidateBackend(t *testing.T) {
	validS3 := &Backend{
		Kind:    util.StrPtr("s3"),
		Bucket:  util.StrPtr("bucket"),
		Profile: util.StrPtr("profile"),
		Region:  util.StrPtr("region"),
	}

	invalidS3 := &Backend{
		Kind: util.StrPtr("s3"),
	}

	validRemote := &Backend{
		Kind:         util.StrPtr("remote"),
		HostName:     util.StrPtr("example.com"),
		Organization: util.StrPtr("org"),
	}

	invalidRemote := &Backend{
		Kind: util.StrPtr("remote"),
	}

	tests := []struct {
		name    string
		backend *Backend
		wantErr bool
	}{
		{"valid-s3", validS3, false},
		{"invalid-s3", invalidS3, true},
		{"valid-remote", validRemote, false},
		{"invalid-remote", invalidRemote, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateBackend(tt.backend, tt.name); (err != nil) != tt.wantErr {
				t.Errorf("ValidateBackend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
