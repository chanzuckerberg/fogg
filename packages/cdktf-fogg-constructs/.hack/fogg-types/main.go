package main

import (
	"fmt"

	// "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/plan"
	ts "github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func main() {
	t := ts.New().
		WithInterface(true).
		WithCustomJsonTag("yaml").
		WithBackupDir("").
		// 	ManageType(v2.Tools{}, ts.TypeOptions{TSType: "{[key: string]: any}"}).
		Add(plan.Component{})

	err := t.ConvertToFile("../../src/imports/fogg-types.generated.ts")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("OK")
}
