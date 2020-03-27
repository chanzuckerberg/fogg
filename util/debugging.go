package util

import (
	"github.com/davecgh/go-spew/spew"
)

// Dump will pretty print whatever is in foo
func Dump(foo interface{}) {
	spew.Dump(foo)
}
