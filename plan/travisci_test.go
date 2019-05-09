package plan

import (
	"testing"

	"github.com/chanzuckerberg/fogg/config/v1"
	"github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/stretchr/testify/assert"
)

var id1, id2 int64

func init() {
	id1 = int64(123456789)
	id1 = int64(987654321)
}

func Test_buildTravisCI_Disabled(t *testing.T) {
	a := assert.New(t)
	{
		c := &v2.Config{
			Tools: v2.Tools{
				TravisCI: &v1.TravisCI{
					Enabled: false,
				},
			},
		}
		p := &Plan{}
		p.Accounts = p.buildAccounts(c)
		tr := p.buildTravisCI(c, "0.1.0")
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
						AccountID: util.Intptr(123),
						Region:    util.StrPtr("us-west-2"),
						Profile:   util.StrPtr("foo"),
						Version:   util.StrPtr("0.12.0"),
					},
				},
				Backend: &v2.Backend{
					Bucket:  util.StrPtr("bucket"),
					Region:  util.StrPtr("us-west-2"),
					Profile: util.StrPtr("profile"),
				},
			},
		},
		Accounts: map[string]v2.Account{
			"foo": v2.Account{
				Common: v2.Common{Providers: &v2.Providers{AWS: &v2.AWSProvider{AccountID: &id1}}},
			},
		},
		Tools: v2.Tools{TravisCI: &v1.TravisCI{
			Enabled:        true,
			AWSIAMRoleName: "rollin",
		}},
	}

	err := c.Validate()
	a.NoError(err)

	p := &Plan{}
	p.Accounts = p.buildAccounts(c)
	tr := p.buildTravisCI(c, "0.1.0")
	a.Len(tr.AWSProfiles, 1)
	a.Equal("foo", tr.AWSProfiles[0].Name)
	a.Equal(id1, tr.AWSProfiles[0].ID)
	a.Equal("rollin", tr.AWSProfiles[0].Role)
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
						AccountID: util.Intptr(123),
						Region:    util.StrPtr("us-west-2"),
						Profile:   util.StrPtr("foo"),
						Version:   util.StrPtr("0.12.0"),
					},
				},
				Backend: &v2.Backend{
					Bucket:  util.StrPtr("bucket"),
					Region:  util.StrPtr("us-west-2"),
					Profile: util.StrPtr("profile"),
				},
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
		Tools: v2.Tools{TravisCI: &v1.TravisCI{
			Enabled:        true,
			AWSIAMRoleName: "rollin",
		}},
	}

	err := c.Validate()
	a.NoError(err)

	p := &Plan{}
	p.Accounts = p.buildAccounts(c)
	tr := p.buildTravisCI(c, "0.1.0")
	a.Len(tr.TestBuckets, 1)
	// 3 because there is always a global
	a.Len(tr.TestBuckets[0], 3)
}
