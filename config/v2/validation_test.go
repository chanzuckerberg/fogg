package v2

import (
	"testing"

	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_nonEmptyString(t *testing.T) {
	r := require.New(t)

	empty := ""
	nonEmpty := "foo"
	r.True(nonEmptyString(&nonEmpty))
	r.False(nonEmptyString(&empty))
	r.False(nonEmptyString(nil))
}

func TestValidateOwnersAccount(t *testing.T) {
	// this will serve as a test for all the fuctions that use validateInheritedStringField, since they are equivalent

	r := require.New(t)
	foo := "foo@example.com"

	// acct owner

	c := confAcctOwner(foo, foo)

	// Both defaults and acct are set
	r.Nil(c.validateInheritedStringField("owner", OwnerGetter, nonEmptyString).ErrorOrNil())

	// defaults unset, still valid
	c = confAcctOwner("", foo)
	r.NoError(c.validateInheritedStringField("owner", OwnerGetter, nonEmptyString).ErrorOrNil())

	// both unset, no longer valid
	c = confAcctOwner("", "")
	r.Equal(2, c.validateInheritedStringField("owner", OwnerGetter, nonEmptyString).Len())
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
			r := require.New(t)
			c := confComponentOwner(tt.def, tt.env, tt.comp)
			e := c.validateInheritedStringField("owner", OwnerGetter, nonEmptyString)
			if tt.errNil {
				r.NoError(e.ErrorOrNil())
			} else {
				r.NotNil(e)
				r.Equal(tt.errz, e.Len())
			}
		})
	}
}

func TestValidateBackend(t *testing.T) {
	r := require.New(t)

	c := confBackendType()

	_, err := c.Validate()
	r.Error(err)
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

func confBackendType() Config {
	return Config{
		Version: 2,
		Defaults: Defaults{
			Common{
				Owner: util.StrPtr("foo@example.com"),
				Backend: &Backend{
					Type:    util.StrPtr("invalid"),
					Bucket:  util.StrPtr("foo"),
					Region:  util.StrPtr("foo"),
					Profile: util.StrPtr("foo"),
				},
				Project:          util.StrPtr("foo"),
				TerraformVersion: util.StrPtr("1.1.1"),
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
	r := require.New(t)
	def := []string{"foo"}
	override := []string{"bar"}

	result := ResolveStringArray(def, override)
	r.Len(result, 1)
	r.Equal("bar", result[0])

	override = nil

	result2 := ResolveStringArray(def, override)
	r.Len(result2, 1)
	r.Equal("foo", result2[0])
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
