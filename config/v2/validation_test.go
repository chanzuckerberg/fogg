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
			fs, _, e := util.TestFs()
			r.NoError(e)
			c := confBackendKind(t, tt.kind)
			_, err := c.Validate(fs)
			if tt.wantErr {
				r.Error(err)
			} else {
				r.NoError(err)
			}
		})
	}
}

func TestValidateGridEndpointRequired(t *testing.T) {
	req := require.New(t)
	enabled := true
	c := &Config{
		Defaults: Defaults{
			Common: Common{
				Grid: &GridConfig{Enabled: &enabled},
			},
		},
	}
	err := c.ValidateGrid()
	req.Error(err)
}

func TestValidateGridEndpointProvided(t *testing.T) {
	req := require.New(t)
	enabled := true
	endpoint := "https://example"
	c := &Config{
		Defaults: Defaults{
			Common: Common{
				Grid: &GridConfig{Enabled: &enabled, Endpoint: &endpoint},
			},
		},
	}
	req.NoError(c.ValidateGrid())
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
				Common: Common{
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

func TestValidateGenericProvider(t *testing.T) {
	validGeneric := &GenericProvider{
		Source: "foo",
		CommonProvider: CommonProvider{
			Version: util.StrPtr("1.1.1"),
		},
	}
	validAssumeRole := &GenericProvider{
		Source: "foo",
		CommonProvider: CommonProvider{
			Version: util.StrPtr("1.1.1"),
		},
		Config: map[string]any{
			"assume_role": map[string]any{
				"role": "foo",
			},
		},
	}

	invalidNothing := &GenericProvider{}
	invalidSource := &GenericProvider{
		Source: "",
	}
	invalidAssumeRoleNothing := &GenericProvider{
		Source: "foo",
		CommonProvider: CommonProvider{
			Version: util.StrPtr("1.1.1"),
		},
		Config: map[string]any{
			"assume_role": map[string]string{},
		},
	}
	invalidAssumeRoleConfig := &GenericProvider{
		Source: "foo",
		CommonProvider: CommonProvider{
			Version: util.StrPtr("1.1.1"),
		},
		Config: map[string]any{
			"assume_role": map[string]string{
				"foo": "bar",
			},
		},
	}

	type args struct {
		p         *GenericProvider
		component string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid generic", args{validGeneric, "valid"}, false},
		{"invalid generic nothing", args{invalidNothing, "invalid"}, true},
		{"invalid generic source", args{invalidSource, "invalid"}, true},
		{"valid generic assume_role", args{validAssumeRole, "valid-role"}, false},
		{"invalid generic assume_role nothing", args{invalidAssumeRoleNothing, "invalid-nothing"}, true},
		{"invalid generic assume_role missing role", args{invalidAssumeRoleConfig, "invalid-role-missing"}, true},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			assert := require.New(t)
			err := tt.args.p.Validate("generic", tt.args.component)

			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
			// t.Errorf("genericProvider.Validate() error = %v, wantErr %v", err, tt.wantErr)
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

func TestValidateFileDependencies(t *testing.T) {
	r := require.New(t)
	fs := afero.NewMemMapFs()

	// Create test files
	err := afero.WriteFile(fs, "file1.txt", []byte("test"), 0644)
	r.NoError(err)
	err = afero.WriteFile(fs, "file2.txt", []byte("test"), 0644)
	r.NoError(err)
	err = afero.WriteFile(fs, "fileTest.txt", []byte("test"), 0644)
	r.NoError(err)
	err = afero.WriteFile(fs, "file-test.txt", []byte("test"), 0644)
	r.NoError(err)

	var cases = []struct {
		label   string
		config  *Config
		wantErr bool
	}{
		{
			label: "valid config",
			config: &Config{
				Defaults: Defaults{
					Common: Common{
						DependsOn: &DependsOn{
							Files: []string{"file1.txt", "file2.txt"},
						},
					},
				},
				Envs: map[string]Env{
					"dev": {
						Common: Common{
							DependsOn: &DependsOn{
								Files: []string{"file1.txt"},
							},
						},
						Components: map[string]Component{
							"web": {
								Common: Common{
									DependsOn: &DependsOn{
										Files: []string{"file2.txt"},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			label: "invalid config",
			config: &Config{
				Envs: map[string]Env{
					"dev": {
						Components: map[string]Component{
							"web": {
								Common: Common{
									DependsOn: &DependsOn{
										Files: []string{"fileTest.txt", "file-test.txt"},
									},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, test := range cases {
		tt := test
		t.Run(tt.label, func(t *testing.T) {
			if err := tt.config.ValidateFileDependencies(fs); (err != nil) != tt.wantErr {
				t.Errorf("Config.ValidateFileDependencies(fs) error = %v, wantErr %v (err != nil) %v", err, tt.wantErr, (err != nil))
			}
		})
	}
}
