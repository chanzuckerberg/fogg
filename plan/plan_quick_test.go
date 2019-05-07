package plan_test

import (
	"fmt"
	"runtime/debug"
	"testing"
	"testing/quick"

	"github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/plan"
)

func add(x, y int) int {
	return x + y
}

func TestValidConfigNoPanic(t *testing.T) {

	// return false if valid + panic
	assertion := func(conf *v2.Config) bool {
		// fmt.Printf("GOT %#v\n\n", pretty.Sprint(conf))
		// validate our configuration
		err := conf.Validate()

		// if config is valid, there should be no panic
		if err == nil {
			fmt.Println("valid")
			defer func() bool {
				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
					debug.PrintStack()
					return false
				}
				return true
			}()

			plan.Eval(conf)

		} else {
			fmt.Println("invalid")
			fmt.Printf("err %s\n", err)
		}

		// config isn't valid so we don't care if we panic or not
		return true

	}
	if err := quick.Check(assertion, &quick.Config{MaxCount: 100}); err != nil {
		t.Error(err)
	}
}
