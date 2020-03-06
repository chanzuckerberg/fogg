package v2_test

import (
	"testing"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/stretchr/testify/assert"
)

func TestResolveTfLint(test *testing.T) {
	a := assert.New(test)
	t := true
	f := false

	data := []struct {
		def    *bool
		over   *bool
		output *bool
	}{
		{nil, nil, &f},
		{nil, &t, &t},
		{nil, &f, &f},
		{&t, nil, &t},
		{&t, &t, &t},
		{&t, &f, &f},
		{&f, nil, &f},
		{&f, &t, &t},
		{&f, &f, &f},
	}
	for _, r := range data {
		test.Run("", func(t *testing.T) {
			def := v2.Common{Tools: &v2.Tools{TfLint: &v2.TfLint{Enabled: r.def}}}
			over := v2.Common{Tools: &v2.Tools{TfLint: &v2.TfLint{Enabled: r.over}}}
			result := v2.ResolveTfLint(def, over)
			a.Equal(r.output, result.Enabled)
		})
	}
}
