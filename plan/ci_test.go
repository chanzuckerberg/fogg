package plan

import (
	"encoding/json"
	"testing"

	v1 "github.com/chanzuckerberg/fogg/config/v1"
	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/stretchr/testify/assert"
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
	a := assert.New(t)
	{
		c := &v2.Config{
			Defaults: v2.Defaults{
				Common: v2.Common{
					Tools: &v2.Tools{
						TravisCI: &v1.TravisCI{
							CommonCI: v1.CommonCI{
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
		a.NotNil(tr)
		a.False(tr.Enabled)
	}
}

func Test_buildTravisCI_Profiles(t *testing.T) {
	a := assert.New(t)

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
					TravisCI: &v1.TravisCI{
						CommonCI: v1.CommonCI{
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
	a.NoError(err)
	a.Len(w, 0)

	p := &Plan{}
	p.Accounts = p.buildAccounts(c)
	tr := p.buildTravisCIConfig(c, "0.1.0")
	a.Len(tr.AWSProfiles, 2)
	a.Contains(tr.AWSProfiles, "profile")
	a.Contains(tr.AWSProfiles, "foo")
	a.Equal(id1.String(), tr.AWSProfiles["foo"].AccountID)
	a.Equal("rollin", tr.AWSProfiles["foo"].RoleName)
}

func Test_buildTravisCI_TestBuckets(t *testing.T) {
	a := assert.New(t)

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
				Tools: &v2.Tools{TravisCI: &v1.TravisCI{
					CommonCI: v1.CommonCI{
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
	a.NoError(err)
	a.Len(w, 0)

	p := &Plan{}
	p.Accounts = p.buildAccounts(c)
	tr := p.buildTravisCIConfig(c, "0.1.0")
	a.NotNil(p.Accounts["foo"].Providers.AWS)
	a.Equal(id1, p.Accounts["foo"].Providers.AWS.AccountID)
	a.Len(tr.TestBuckets, 1)
	a.Len(tr.TestBuckets[0], 2)
}
