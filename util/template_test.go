package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDict(t *testing.T) {
	m := make(map[string]string)
	m["foo"] = "bar"
	r := dict(m)
	require.NotNil(t, r)
	require.IsType(t, map[string]interface{}{}, r)
	require.Equal(t, "bar", r["foo"])
}

func TestRenderHCLBody(t *testing.T) {
	tests := []struct {
		name   string
		config map[string]any
		indent int
		want   string
	}{
		{
			name:   "nil config",
			config: nil,
			indent: 2,
			want:   "",
		},
		{
			name:   "empty config",
			config: map[string]any{},
			indent: 2,
			want:   "",
		},
		{
			name:   "string value",
			config: map[string]any{"region": "us-east-1"},
			indent: 2,
			want:   "  region = \"us-east-1\"\n",
		},
		{
			name:   "bool value",
			config: map[string]any{"enabled": true},
			indent: 2,
			want:   "  enabled = true\n",
		},
		{
			name:   "int value",
			config: map[string]any{"count": float64(42)},
			indent: 2,
			want:   "  count = 42\n",
		},
		{
			name:   "list value",
			config: map[string]any{"files": []any{"a.txt", "b.txt"}},
			indent: 2,
			want:   "  files = [\"a.txt\", \"b.txt\"]\n",
		},
		{
			name: "nested block",
			config: map[string]any{
				"assume_role": map[string]any{
					"role_arn": "arn:aws:iam::123:role/foo",
				},
			},
			indent: 2,
			want:   "  assume_role {\n    role_arn = \"arn:aws:iam::123:role/foo\"\n  }\n",
		},
		{
			name: "sorted keys",
			config: map[string]any{
				"region":      "us-west-2",
				"assume_role": map[string]any{"role_arn": "arn"},
			},
			indent: 2,
			want:   "  assume_role {\n    role_arn = \"arn\"\n  }\n  region = \"us-west-2\"\n",
		},
		{
			name:   "custom indent",
			config: map[string]any{"key": "val"},
			indent: 4,
			want:   "    key = \"val\"\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderHCLBody(tt.config, tt.indent)
			require.Equal(t, tt.want, got)
		})
	}
}
