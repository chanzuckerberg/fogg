package v2

import (
	"testing"

	"github.com/chanzuckerberg/fogg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/spf13/afero"
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
		{"v2_full", false},
		{"v2_minimal_valid", false},
		{"v2_invalid_aws_provider", true},
	}
	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			r := require.New(t)
			fs, _, _ := util.TestFs()

			b, e := util.TestFile(tt.fileName)
			r.NoError(e)
			e = afero.WriteFile(fs, "fogg.json", b, 0644)
			r.NoError(e)

			c, e := ReadConfig(b, fs, "fogg.json")
			r.NoError(e)
			r.NotNil(c)

			if err := c.ValidateAWSProviders(); (err != nil) != tt.wantErr {
				t.Errorf("Config.ValidateAWSProviders() error = %v, wantErr %v (err != nil) %v", err, tt.wantErr, (err != nil))
			}
		})
	}
}
