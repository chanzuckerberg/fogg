package plan

import (
	"testing"

	"github.com/chanzuckerberg/fogg/config"
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
		c := &config.Config{
			TravisCI: &config.TravisCI{
				Enabled: false,
			},
		}
		p := &Plan{}
		tr := p.buildTravisCI(c)
		a.NotNil(tr)
		a.False(tr.Enabled)
	}
}
func Test_buildTravisCI_Profiles(t *testing.T) {
	a := assert.New(t)

	c := &config.Config{
		Accounts: map[string]config.Account{
			"foo": config.Account{
				AccountID: &id1,
			},
		},
		TravisCI: &config.TravisCI{
			Enabled:        true,
			AWSIAMRoleName: "rollin",
			HubAccountName: "hub",
		},
	}
	p := &Plan{}
	tr := p.buildTravisCI(c)
	a.Len(tr.AWSProfiles, 1)
	a.Equal(tr.AWSProfiles[0].Name, "foo")
	a.Equal(tr.AWSProfiles[0].ID, id1)
	a.Equal(tr.AWSProfiles[0].Role, "rollin")
	a.Equal(tr.AWSProfiles[0].HubAccountName, "hub")
}

func Test_buildTravisCI_TestShards(t *testing.T) {
	a := assert.New(t)

	c := &config.Config{
		Accounts: map[string]config.Account{
			"foo": config.Account{
				AccountID: &id1,
			},
			"bar": config.Account{
				AccountID: &id2,
			},
		},
		TravisCI: &config.TravisCI{
			Enabled:        true,
			AWSIAMRoleName: "rollin",
			HubAccountName: "hub",
		},
	}

	p := &Plan{}
	tr := p.buildTravisCI(c)
	a.Len(tr.TestShards, 1)
	a.Len(tr.TestShards[0], 2)
}
