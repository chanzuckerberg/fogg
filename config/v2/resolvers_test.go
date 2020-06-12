package v2_test

import (
	"fmt"
	"testing"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/stretchr/testify/require"
)

func TestResolveTfLint(t *testing.T) {
	r := require.New(t)
	tru := true
	fal := false
	optEmpty := ""
	optSet := "--deep"
	optOver := "--deep"

	data := []struct {
		enabledDef  *bool
		optionDef   *string
		enabledOver *bool
		optionOver  *string
		enabledOut  *bool
		optionOut   *string
	}{
		{nil, nil, &fal, nil, &fal, &optEmpty},
		{nil, nil, &tru, &optEmpty, &tru, &optEmpty},
		{nil, nil, &fal, &optSet, &fal, &optSet},
		{&tru, &optEmpty, nil, nil, &tru, &optEmpty},
		{&tru, &optEmpty, &tru, &optEmpty, &tru, &optEmpty},
		{&tru, &optEmpty, &fal, &optSet, &fal, &optSet},
		{&fal, &optSet, nil, nil, &fal, &optSet},
		{&fal, &optSet, &tru, &optEmpty, &tru, &optSet},
		{&fal, &optSet, &fal, &optOver, &fal, &optOver},
	}
	for _, tt := range data {
		t.Run("", func(t *testing.T) {
			fmt.Println("boo")
			def := v2.Common{Tools: &v2.Tools{TfLint: &v2.TfLint{Enabled: tt.enabledDef, Options: tt.optionDef}}}
			over := v2.Common{Tools: &v2.Tools{TfLint: &v2.TfLint{Enabled: tt.enabledOver, Options: tt.optionOver}}}
			result := v2.ResolveTfLint(def, over)
			r.Equal(tt.enabledOut, result.Enabled)
			r.Equal(tt.optionOut, result.Options)
		})
	}
}
