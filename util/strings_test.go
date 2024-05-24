package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertCamelToSnake(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Single word", "hello", "hello"},
		{"All lowercase", "helloWorld", "hello_world"},
		{"Mixed case", "helloWorldFooBar", "hello_world_foo_bar"},
		{"All uppercase", "HELLO", "hello"},
		{"Mixed case with hyphens", "hello-worldFoo-bar", "hello_world_foo_bar"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := require.New(t)
			result := ConvertToSnake(test.input)
			r.Equal(test.expected, result)
		})
	}
}
