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
		p.Accounts = p.buildAccounts(c)
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
				Project:          util.StrPtr("foo"),
				Owner:            util.StrPtr("bar"),
				TerraformVersion: util.StrPtr("0.1.0"),
				Providers: &v2.Providers{
					AWS: &v2.AWSProvider{
						AccountID: util.JsonNumberPtr(123),
						Region:    util.StrPtr("us-west-2"),
						Profile:   util.StrPtr("foo"),
						Version:   util.StrPtr("0.12.0"),
					},
				},
				Backend: &v2.Backend{
					Bucket:    util.StrPtr("bucket"),
					Region:    util.StrPtr("us-west-2"),
					Profile:   util.StrPtr("profile"),
					AccountID: util.StrPtr("some account id"),
				},
				Tools: &v2.Tools{
					TravisCI: &v2.TravisCI{
						CommonCI: v2.CommonCI{
							Enabled:        &tr,
							AWSIAMRoleName: util.StrPtr("rollin"),
						},
					}},
			},
		},
		Accounts: map[string]v2.Account{
			"foo": v2.Account{
				Common: v2.Common{Providers: &v2.Providers{AWS: &v2.AWSProvider{AccountID: &id1}}},
			},
		},
	}

	w, err := c.Validate()
	r.NoError(err)
	r.Len(w, 0)

	p := &Plan{}
	p.Accounts = p.buildAccounts(c)
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
				Project:          util.StrPtr("foo"),
				Owner:            util.StrPtr("bar"),
				TerraformVersion: util.StrPtr("0.1.0"),
				Providers: &v2.Providers{
					AWS: &v2.AWSProvider{
						AccountID: util.JsonNumberPtr(123),
						Region:    util.StrPtr("us-west-2"),
						Profile:   util.StrPtr("foo"),
						Version:   util.StrPtr("0.12.0"),
					},
				},
				Backend: &v2.Backend{
					Bucket:    util.StrPtr("bucket"),
					Region:    util.StrPtr("us-west-2"),
					Profile:   util.StrPtr("profile"),
					AccountID: util.StrPtr("some account id"),
				},
				Tools: &v2.Tools{TravisCI: &v2.TravisCI{
					CommonCI: v2.CommonCI{
						Enabled:        &tr,
						AWSIAMRoleName: util.StrPtr("rollin"),
					},
				}},
			},
		},
		Accounts: map[string]v2.Account{
			"foo": v2.Account{
				Common: v2.Common{Providers: &v2.Providers{AWS: &v2.AWSProvider{AccountID: &id1}}},
			},
			"bar": v2.Account{
				Common: v2.Common{Providers: &v2.Providers{AWS: &v2.AWSProvider{AccountID: &id2}}},
			},
		},
	}

	w, err := c.Validate()
	r.NoError(err)
	r.Len(w, 0)

	p := &Plan{}
	p.Accounts = p.buildAccounts(c)
	tr := p.buildTravisCIConfig(c, "0.1.0")
	r.NotNil(p.Accounts["foo"].Providers.AWS)
	r.Equal(id1, p.Accounts["foo"].Providers.AWS.AccountID)
	r.Len(tr.TestBuckets, 1)
	r.Len(tr.TestBuckets[0], 2)
}

func Test_buildCircleCI_Profiles(t *testing.T) {
	r := require.New(t)

	c := &v2.Config{
		Version: 2,
		Defaults: v2.Defaults{
			Common: v2.Common{
				Project:          util.StrPtr("foo"),
				Owner:            util.StrPtr("bar"),
				TerraformVersion: util.StrPtr("0.1.0"),
				Providers: &v2.Providers{
					AWS: &v2.AWSProvider{
						AccountID: util.JsonNumberPtr(123),
						Region:    util.StrPtr("us-west-2"),
						Profile:   util.StrPtr("foo"),
						Version:   util.StrPtr("0.12.0"),
					},
				},
				Backend: &v2.Backend{
					Bucket:    util.StrPtr("bucket"),
					Region:    util.StrPtr("us-west-2"),
					Profile:   util.StrPtr("profile"),
					AccountID: util.StrPtr("some account id"),
				},
				Tools: &v2.Tools{
					CircleCI: &v2.CircleCI{
						CommonCI: v2.CommonCI{
							Enabled:        &tr,
							AWSIAMRoleName: util.StrPtr("rollin"),
						},
					}},
			},
		},
		Accounts: map[string]v2.Account{
			"foo": v2.Account{
				Common: v2.Common{Providers: &v2.Providers{AWS: &v2.AWSProvider{AccountID: &id1}}},
			},
		},
	}

	w, err := c.Validate()
	r.NoError(err)
	r.Len(w, 0)

	p := &Plan{}
	p.Accounts = p.buildAccounts(c)
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
				Project:          util.StrPtr("foo"),
				Owner:            util.StrPtr("bar"),
				TerraformVersion: util.StrPtr("0.1.0"),
				Providers: &v2.Providers{
					AWS: &v2.AWSProvider{
						AccountID: util.JsonNumberPtr(123),
						Region:    util.StrPtr("us-west-2"),
						Profile:   util.StrPtr("foo"),
						Version:   util.StrPtr("0.12.0"),
					},
				},
				Backend: &v2.Backend{
					Bucket:    util.StrPtr("bucket"),
					Region:    util.StrPtr("us-west-2"),
					Profile:   util.StrPtr("profile"),
					AccountID: util.StrPtr("some account id"),
				},
				Tools: &v2.Tools{
					CircleCI: &v2.CircleCI{
						CommonCI: v2.CommonCI{
							Enabled:        &tr,
							AWSIAMRoleName: util.StrPtr("rollin"),
							Providers: map[string]v2.CIProviderConfig{
								"aws": v2.CIProviderConfig{
									Disabled: true,
								},
							},
						},
					}},
			},
		},
		Accounts: map[string]v2.Account{
			"foo": v2.Account{
				Common: v2.Common{Providers: &v2.Providers{AWS: &v2.AWSProvider{AccountID: &id1}}},
			},
		},
	}

	w, err := c.Validate()
	r.NoError(err)
	r.Len(w, 0)

	p := &Plan{}
	p.Accounts = p.buildAccounts(c)
	circle := p.buildCircleCIConfig(c, "0.1.0")
	r.Len(circle.AWSProfiles, 0)
}
