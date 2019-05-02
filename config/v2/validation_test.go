package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_nonEmptyString(t *testing.T) {
	a := assert.New(t)

	a.True(nonEmptyString("foo"))
	a.False(nonEmptyString(""))
}

func TestValidateOwners(t *testing.T) {
	// this will serve as a test for all the fuctions that use validateInheritedStringField, since they are equivalent

	a := assert.New(t)
	foo := "foo@example.com"

	// acct owner
	{
		c := confAcctOwner(foo, foo)

		// Both defaults and acct are set
		a.Nil(c.validateOwners())

		// defaults unset, still valid
		c = confAcctOwner("", foo)
		a.Nil(c.validateOwners())

		// both unset, no longer valid
		c = confAcctOwner("", "")
		a.Equal(1, c.validateOwners().Len())
	}

	// component owner
	{
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
				c := confComponentOwner(tt.def, tt.env, tt.comp)
				e := c.validateOwners()
				if tt.errNil {
					a.Nil(e)
				} else {
					a.NotNil(e)
					a.Equal(tt.errz, e.Len())
				}
			})
		}
	}
}

func confAcctOwner(def, acct string) Config {
	return Config{
		Defaults: Defaults{
			common{
				Owner: def,
			},
		},
		Accounts: map[string]Account{
			"foo": Account{
				common{
					Owner: acct,
				},
			},
		},
	}
}

func confComponentOwner(def, env, component string) Config {
	return Config{
		Defaults: Defaults{
			common{
				Owner: def,
			},
		},
		Envs: map[string]Env{
			"bar": Env{
				common: common{
					Owner: env,
				},
				Components: map[string]Component{
					"bam": Component{
						common: common{
							Owner: component,
						},
					},
				},
			},
		},
	}
}
