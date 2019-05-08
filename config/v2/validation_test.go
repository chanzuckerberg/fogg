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

func TestValidateOwnersAccount(t *testing.T) {
	// this will serve as a test for all the fuctions that use validateInheritedStringField, since they are equivalent

	a := assert.New(t)
	foo := "foo@example.com"

	// acct owner

	c := confAcctOwner(foo, foo)

	// Both defaults and acct are set
	a.Nil(c.validateOwners().ErrorOrNil())

	// defaults unset, still valid
	c = confAcctOwner("", foo)
	a.NoError(c.validateOwners().ErrorOrNil())

	// both unset, no longer valid
	c = confAcctOwner("", "")
	a.Equal(2, c.validateOwners().Len())
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
			e := c.validateOwners()
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
				Owner: def,
				Backend: Backend{
					Bucket:  "foo",
					Region:  "foo",
					Profile: "foo",
				},
			},
		},
		Accounts: map[string]Account{
			"foo": Account{
				Common{
					Owner: acct,
				},
			},
		},
		Global: Component{
			Common: Common{
				Owner: acct,
			},
		},
	}
}

func confComponentOwner(def, env, component string) Config {
	return Config{
		Defaults: Defaults{
			Common{
				Owner: def,
				Backend: Backend{
					Bucket:  "foo",
					Region:  "foo",
					Profile: "foo",
				},
			},
		},
		Envs: map[string]Env{
			"bar": Env{
				Common: Common{
					Owner: env,
				},
				Components: map[string]Component{
					"bam": Component{
						Common: Common{
							Owner: component,
						},
					},
				},
			},
		},
		Global: Component{
			Common: Common{
				Owner: "foo",
			},
		},
	}
}
