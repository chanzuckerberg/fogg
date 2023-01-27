package v2

import (
	"testing"

	"github.com/chanzuckerberg/fogg/util"
	"github.com/jinzhu/copier"
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

	for _, test := range cases {
		tt := test
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

func TestValidateBackends(t *testing.T) {
	var cases = []struct {
		kind    string
		wantErr bool
	}{
		{"invalid", true},
		{"s3", false},
		{"remote", false},
	}

	for _, test := range cases {
		tt := test
		t.Run(tt.kind, func(t *testing.T) {
			r := require.New(t)

			c := confBackendKind(t, tt.kind)
			_, err := c.Validate()
			if tt.wantErr {
				r.Error(err)
			} else {
				r.NoError(err)
			}
		})
	}
}

func confBackendKind(t *testing.T, kind string) Config {
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
			"foo": {
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
			"bar": {
				Common: Common{
					Owner: util.StrPtr(env),
				},
				Components: map[string]Component{
					"bam": {
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

func TestValidateAWSProvider(t *testing.T) {
	validProfile := &AWSProvider{
		AccountID: util.JSONNumberPtr(123456),
		Profile:   util.StrPtr("my-profile"),
		Region:    util.StrPtr("us-sw-12"),
		CommonProvider: CommonProvider{
			Version: util.StrPtr("1.1.1"),
		},
	}

	invalidNothing := &AWSProvider{}
	invalidBoth := &AWSProvider{
		AccountID: util.JSONNumberPtr(123456),
		Profile:   util.StrPtr("my-profile-name"),
		Role:      util.StrPtr("my-role-name"),
		Region:    util.StrPtr("us-sw-12"),
		CommonProvider: CommonProvider{
			Version: util.StrPtr("1.1.1"),
		},
	}

	validRole := &AWSProvider{
		AccountID: util.JSONNumberPtr(123456),
		Role:      util.StrPtr("my-role-name"),
		Region:    util.StrPtr("us-sw-12"),
		CommonProvider: CommonProvider{
			Version: util.StrPtr("1.1.1"),
		},
	}

	type args struct {
		p         *AWSProvider
		component string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid profile", args{validProfile, "valid"}, false},
		{"invalid", args{invalidNothing, "invalid"}, true},
		{"valid role", args{validRole, "valid-role"}, false},
		{"invalid both", args{invalidBoth, "invalid-both"}, true},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateAWSProvider(tt.args.p, tt.args.component); (err != nil) != tt.wantErr {
				t.Errorf("ValidateAWSProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
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
	for _, test := range tests {
		tt := test
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
	for _, test := range tests {
		tt := test
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
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateBackend(tt.backend, tt.name); (err != nil) != tt.wantErr {
				t.Errorf("ValidateBackend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
