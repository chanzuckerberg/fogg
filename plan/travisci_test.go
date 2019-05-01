package plan

import (
	"testing"

	"github.com/chanzuckerberg/fogg/config/v1"
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
		c := &v1.Config{
			TravisCI: &v1.TravisCI{
				Enabled: false,
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

	c := &v1.Config{
		Accounts: map[string]v1.Account{
			"foo": v1.Account{
				AccountID: &id1,
			},
		},
		TravisCI: &v1.TravisCI{
			Enabled:        true,
			AWSIAMRoleName: "rollin",
		},
	}
	p := &Plan{}
	p.Accounts = p.buildAccounts(c)
	tr := p.buildTravisCI(c, "0.1.0")
	a.Len(tr.AWSProfiles, 1)
	a.Equal(tr.AWSProfiles[0].Name, "foo")
	a.Equal(tr.AWSProfiles[0].ID, id1)
	a.Equal(tr.AWSProfiles[0].Role, "rollin")
}

func Test_buildTravisCI_TestBuckets(t *testing.T) {
	a := assert.New(t)

	c := &v1.Config{
		Accounts: map[string]v1.Account{
			"foo": v1.Account{
				AccountID: &id1,
			},
			"bar": v1.Account{
				AccountID: &id2,
			},
		},
		TravisCI: &v1.TravisCI{
			Enabled:        true,
			AWSIAMRoleName: "rollin",
		},
	}

	p := &Plan{}
	p.Accounts = p.buildAccounts(c)
	tr := p.buildTravisCI(c, "0.1.0")
	a.Len(tr.TestBuckets, 1)
	// 3 because there is always a global
	a.Len(tr.TestBuckets[0], 3)
}
