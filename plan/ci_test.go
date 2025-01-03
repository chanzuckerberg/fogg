package plan

import (
	"encoding/json"
	"testing"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/stretchr/testify/require"
)

var id1, id2 json.Number

var f, tr bool

func init() {
	id1 = json.Number("123456789")
	id2 = json.Number("987654321")

	f = false
	tr = true
}

func Test_buildTravisCI_Disabled(t *testing.T) {
	r := require.New(t)
	{
		c := &v2.Config{
			Defaults: v2.Defaults{
				Common: v2.Common{
					Tools: &v2.Tools{
						TravisCI: &v2.TravisCI{
							CommonCI: v2.CommonCI{
								Enabled: &f,
							},
						},
					},
				},
			},
		}
		p := &Plan{}
		accts := p.buildAccounts(c)
		p.Accounts = accts
		tr := p.buildTravisCIConfig(c, "0.1.0")
		r.NotNil(tr)
		r.False(tr.Enabled)
	}
}

func Test_buildTravisCI_Profiles(t *testing.T) {
	r := require.New(t)

	c := &v2.Config{
		Version: 2,
		Defaults: v2.Defaults{
			Common: v2.Common{
				Project:          util.Ptr("foo"),
				Owner:            util.Ptr("bar"),
				TerraformVersion: util.Ptr("0.1.0"),
				Providers: &v2.Providers{
					AWS: &v2.AWSProvider{
						AccountID: util.JSONNumberPtr(123),
						Region:    util.Ptr("us-west-2"),
						Profile:   util.Ptr("foo"),
						CommonProvider: v2.CommonProvider{
							Version: util.StrPtr("0.12.0"),
						},
					},
				},
				Backend: &v2.Backend{
					Bucket:    util.Ptr("bucket"),
					Region:    util.Ptr("us-west-2"),
					Profile:   util.Ptr("profile"),
					AccountID: util.Ptr("some account id"),
				},
				Tools: &v2.Tools{
					TravisCI: &v2.TravisCI{
						CommonCI: v2.CommonCI{
							Enabled:        &tr,
							AWSIAMRoleName: util.Ptr("rollin"),
						},
					}},
			},
		},
		Accounts: map[string]v2.Account{
			"foo": {
				Common: v2.Common{Providers: &v2.Providers{AWS: &v2.AWSProvider{AccountID: &id1}}},
			},
		},
	}
	fs, _, e := util.TestFs()
	r.NoError(e)
	w, err := c.Validate(fs)
	r.NoError(err)
	r.Len(w, 0)

	p := &Plan{}
	accts := p.buildAccounts(c)
	r.Len(accts, 1)
	p.Accounts = accts
	tr := p.buildTravisCIConfig(c, "0.1.0")
	r.Len(tr.AWSProfiles, 2)
	r.Contains(tr.AWSProfiles, "profile")
	r.Contains(tr.AWSProfiles, "foo")
	r.Equal(id1.String(), tr.AWSProfiles["foo"].AccountID)
	r.Equal("rollin", tr.AWSProfiles["foo"].RoleName)
}

func Test_buildTravisCI_TestBuckets(t *testing.T) {
	r := require.New(t)

	c := &v2.Config{
		Version: 2,
		Defaults: v2.Defaults{
			Common: v2.Common{
				Project:          util.Ptr("foo"),
				Owner:            util.Ptr("bar"),
				TerraformVersion: util.Ptr("0.1.0"),
				Providers: &v2.Providers{
					AWS: &v2.AWSProvider{
						AccountID: util.JSONNumberPtr(123),
						Region:    util.Ptr("us-west-2"),
						Profile:   util.Ptr("foo"),
						CommonProvider: v2.CommonProvider{
							Version: util.StrPtr("0.12.0"),
						},
					},
				},
				Backend: &v2.Backend{
					Bucket:    util.Ptr("bucket"),
					Region:    util.Ptr("us-west-2"),
					Profile:   util.Ptr("profile"),
					AccountID: util.Ptr("some account id"),
				},
				Tools: &v2.Tools{TravisCI: &v2.TravisCI{
					CommonCI: v2.CommonCI{
						Enabled:        &tr,
						AWSIAMRoleName: util.Ptr("rollin"),
					},
				}},
			},
		},
		Accounts: map[string]v2.Account{
			"foo": {
				Common: v2.Common{Providers: &v2.Providers{AWS: &v2.AWSProvider{AccountID: &id1}}},
			},
			"bar": {
				Common: v2.Common{Providers: &v2.Providers{AWS: &v2.AWSProvider{AccountID: &id2}}},
			},
		},
	}
	fs, _, e := util.TestFs()
	r.NoError(e)
	w, err := c.Validate(fs)
	r.NoError(err)
	r.Len(w, 0)

	p := &Plan{}
	accts := p.buildAccounts(c)
	p.Accounts = accts
	tr := p.buildTravisCIConfig(c, "0.1.0")
	r.NotNil(p.Accounts["foo"].ProviderConfiguration.AWS)
	r.Equal(id1, p.Accounts["foo"].ProviderConfiguration.AWS.AccountID)
	r.Len(tr.TestBuckets, 1)
	r.Len(tr.TestBuckets[0], 2)
}

func Test_buildCircleCI_Profiles(t *testing.T) {
	r := require.New(t)

	c := &v2.Config{
		Version: 2,
		Defaults: v2.Defaults{
			Common: v2.Common{
				Project:          util.Ptr("foo"),
				Owner:            util.Ptr("bar"),
				TerraformVersion: util.Ptr("0.1.0"),
				Providers: &v2.Providers{
					AWS: &v2.AWSProvider{
						AccountID: util.JSONNumberPtr(123),
						Region:    util.Ptr("us-west-2"),
						Profile:   util.Ptr("foo"),
						CommonProvider: v2.CommonProvider{
							Version: util.StrPtr("0.12.0"),
						},
					},
				},
				Backend: &v2.Backend{
					Bucket:    util.Ptr("bucket"),
					Region:    util.Ptr("us-west-2"),
					Profile:   util.Ptr("profile"),
					AccountID: util.Ptr("some account id"),
				},
				Tools: &v2.Tools{
					CircleCI: &v2.CircleCI{
						CommonCI: v2.CommonCI{
							Enabled:        &tr,
							AWSIAMRoleName: util.Ptr("rollin"),
						},
					}},
			},
		},
		Accounts: map[string]v2.Account{
			"foo": {
				Common: v2.Common{Providers: &v2.Providers{AWS: &v2.AWSProvider{AccountID: &id1}}},
			},
		},
	}
	fs, _, e := util.TestFs()
	r.NoError(e)
	w, err := c.Validate(fs)
	r.NoError(err)
	r.Len(w, 0)

	p := &Plan{}
	accts := p.buildAccounts(c)
	p.Accounts = accts
	circle := p.buildCircleCIConfig(c, "0.1.0")
	r.Len(circle.AWSProfiles, 2)
	r.Contains(circle.AWSProfiles, "profile")
	r.Contains(circle.AWSProfiles, "foo")
	r.Equal(id1.String(), circle.AWSProfiles["foo"].AccountID)
	r.Equal("rollin", circle.AWSProfiles["foo"].RoleName)
}

func Test_buildCircleCI_ProfilesDisabled(t *testing.T) {
	r := require.New(t)

	c := &v2.Config{
		Version: 2,
		Defaults: v2.Defaults{
			Common: v2.Common{
				Project:          util.Ptr("foo"),
				Owner:            util.Ptr("bar"),
				TerraformVersion: util.Ptr("0.1.0"),
				Providers: &v2.Providers{
					AWS: &v2.AWSProvider{
						AccountID: util.JSONNumberPtr(123),
						Region:    util.Ptr("us-west-2"),
						Profile:   util.Ptr("foo"),
						CommonProvider: v2.CommonProvider{
							Version: util.StrPtr("0.12.0"),
						},
					},
				},
				Backend: &v2.Backend{
					Bucket:    util.Ptr("bucket"),
					Region:    util.Ptr("us-west-2"),
					Profile:   util.Ptr("profile"),
					AccountID: util.Ptr("some account id"),
				},
				Tools: &v2.Tools{
					CircleCI: &v2.CircleCI{
						CommonCI: v2.CommonCI{
							Enabled:        &tr,
							AWSIAMRoleName: util.Ptr("rollin"),
							Providers: map[string]v2.CIProviderConfig{
								"aws": {
									Disabled: true,
								},
							},
						},
					}},
			},
		},
		Accounts: map[string]v2.Account{
			"foo": {
				Common: v2.Common{Providers: &v2.Providers{AWS: &v2.AWSProvider{AccountID: &id1}}},
			},
		},
	}
	fs, _, e := util.TestFs()
	r.NoError(e)
	w, err := c.Validate(fs)
	r.NoError(err)
	r.Len(w, 0)

	p := &Plan{}
	accts := p.buildAccounts(c)
	p.Accounts = accts
	circle := p.buildCircleCIConfig(c, "0.1.0")
	r.Len(circle.AWSProfiles, 0)
}

func Test_parseScopes(t *testing.T) {
	r := require.New(t)

	t.Run("No new scopes", func(t *testing.T) {
		defaultScopes := map[string]jsScope{}
		newScopes := []jsScope{}
		result, script := parseJsScopes(&defaultScopes, newScopes)
		r.NotNil(result)
		r.Equal(noCALoginRequired, script)
		r.Empty(*result)
	})

	t.Run("No new scopes, keep defaults", func(t *testing.T) {
		defaultScopes := map[string]jsScope{
			"@scope1": {Name: "@scope1", RegistryUrl: "https://registry.npmjs.org"},
		}
		newScopes := []jsScope{}
		result, script := parseJsScopes(&defaultScopes, newScopes)
		r.NotNil(result)
		r.Equal(noCALoginRequired, script)
		r.Len(*result, 1)
		r.Equal("https://registry.npmjs.org", (*result)["@scope1"].RegistryUrl)
	})

	t.Run("New scopes without CodeArtifact URLs", func(t *testing.T) {
		defaultScopes := map[string]jsScope{}
		newScopes := []jsScope{
			{Name: "@scope1", RegistryUrl: "https://registry.npmjs.org"},
			{Name: "@scope2", RegistryUrl: "https://registry.yarnpkg.com"},
		}
		result, script := parseJsScopes(&defaultScopes, newScopes)
		r.NotNil(result)
		r.Equal(noCALoginRequired, script)
		r.Len(*result, 2)
		r.Equal("https://registry.npmjs.org", (*result)["@scope1"].RegistryUrl)
		r.Equal("https://registry.yarnpkg.com", (*result)["@scope2"].RegistryUrl)
	})

	t.Run("New scopes without CodeArtifact URLs, merge with defaults", func(t *testing.T) {
		defaultScopes := map[string]jsScope{
			"@scope1": {Name: "@scope1", RegistryUrl: "https://registry.yarnpkg.com"},
		}
		newScopes := []jsScope{
			{Name: "@scope1", RegistryUrl: "https://registry.npmjs.org"},
			{Name: "@scope2", RegistryUrl: "https://registry.yarnpkg.com"},
		}
		result, script := parseJsScopes(&defaultScopes, newScopes)
		r.NotNil(result)
		r.Equal(noCALoginRequired, script)
		r.Len(*result, 2)
		r.Equal("https://registry.npmjs.org", (*result)["@scope1"].RegistryUrl)
		r.Equal("https://registry.yarnpkg.com", (*result)["@scope2"].RegistryUrl)
	})

	t.Run("New scopes with CodeArtifact URLs", func(t *testing.T) {
		mergedScopes := map[string]jsScope{}
		newScopes := []jsScope{
			{Name: "@scope1", RegistryUrl: "https://registry.npmjs.org"},
			{Name: "@scope2", RegistryUrl: "https://domain-123456789012.d.codeartifact.us-west-2.amazonaws.com/npm/repo-name/"},
			{Name: "@scope3", RegistryUrl: "https://domain-123456789012.d.codeartifact.us-west-2.amazonaws.com/npm/repo-name/"},
			{Name: "@scope4", RegistryUrl: "https://domain-123456789012.d.codeartifact.us-east-1.amazonaws.com/npm/repo-name/"},
			{Name: "@scope5", RegistryUrl: "https://domain-210987654321.d.codeartifact.us-west-2.amazonaws.com/npm/repo-name/"},
		}
		result, script := parseJsScopes(&mergedScopes, newScopes)
		r.NotNil(result)
		r.NotNil(script)
		r.Len(*result, 5)
		r.Equal("https://registry.npmjs.org", (*result)["@scope1"].RegistryUrl)
		r.Equal("https://domain-123456789012.d.codeartifact.us-west-2.amazonaws.com/npm/repo-name/", (*result)["@scope2"].RegistryUrl)
		r.Equal("https://domain-123456789012.d.codeartifact.us-west-2.amazonaws.com/npm/repo-name/", (*result)["@scope3"].RegistryUrl)
		r.Equal("https://domain-123456789012.d.codeartifact.us-east-1.amazonaws.com/npm/repo-name/", (*result)["@scope4"].RegistryUrl)
		r.Equal("https://domain-210987654321.d.codeartifact.us-west-2.amazonaws.com/npm/repo-name/", (*result)["@scope5"].RegistryUrl)
		expectedScript := `aws codeartifact login --tool npm --repository repo-name --domain domain --domain-owner 123456789012 --region us-east-1 --namespace "@scope4";` +
			`aws codeartifact login --tool npm --repository repo-name --domain domain --domain-owner 123456789012 --region us-west-2 --namespace "@scope2";` +
			`aws codeartifact login --tool npm --repository repo-name --domain domain --domain-owner 123456789012 --region us-west-2 --namespace "@scope3";` +
			`aws codeartifact login --tool npm --repository repo-name --domain domain --domain-owner 210987654321 --region us-west-2 --namespace "@scope5"`
		r.Equal(expectedScript, script)
	})

	t.Run("Merge existing scopes with new scopes", func(t *testing.T) {
		mergedScopes := map[string]jsScope{
			"@scope1": {Name: "@scope1", RegistryUrl: "https://registry.npmjs.org"},
		}
		newScopes := []jsScope{
			{Name: "@scope2", RegistryUrl: "https://domain-123456789012.d.codeartifact.us-west-2.amazonaws.com/npm/repo-name/"},
		}
		result, script := parseJsScopes(&mergedScopes, newScopes)
		r.NotNil(result)
		r.NotNil(script)
		r.Len(*result, 2)
		r.Equal("https://registry.npmjs.org", (*result)["@scope1"].RegistryUrl)
		r.Equal("https://domain-123456789012.d.codeartifact.us-west-2.amazonaws.com/npm/repo-name/", (*result)["@scope2"].RegistryUrl)
		r.Contains(script, `aws codeartifact login --tool npm --repository repo-name --domain domain --domain-owner 123456789012 --region us-west-2 --namespace "@scope2"`)
	})
}
