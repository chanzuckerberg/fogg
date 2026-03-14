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
		name    string
		config  map[string]any
		indent  int
		objects []string
		want    string
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
			name: "nested block default",
			config: map[string]any{
				"assume_role": map[string]any{
					"role_arn": "arn:aws:iam::123:role/foo",
				},
			},
			indent: 2,
			want:   "  assume_role {\n    role_arn = \"arn:aws:iam::123:role/foo\"\n  }\n",
		},
		{
			name: "nested attribute object via objects list",
			config: map[string]any{
				"assume_role": map[string]any{
					"role_arn": "arn:aws:iam::123:role/foo",
				},
			},
			indent:  2,
			objects: []string{"assume_role"},
			want:    "  assume_role = {\n    role_arn = \"arn:aws:iam::123:role/foo\"\n  }\n",
		},
		{
			name:   "empty block default",
			config: map[string]any{"features": map[string]any{}},
			indent: 2,
			want:   "  features {}\n",
		},
		{
			name:    "empty attribute object via objects list",
			config:  map[string]any{"features": map[string]any{}},
			indent:  2,
			objects: []string{"features"},
			want:    "  features = {}\n",
		},
		{
			name: "deep nested blocks default",
			config: map[string]any{
				"update": map[string]any{
					"server": "192.168.0.1",
					"gssapi": map[string]any{
						"realm": "EXAMPLE.COM",
					},
				},
			},
			indent: 2,
			want:   "  update {\n    gssapi {\n      realm = \"EXAMPLE.COM\"\n    }\n    server = \"192.168.0.1\"\n  }\n",
		},
		{
			name: "mixed block and object",
			config: map[string]any{
				"region":      "us-west-2",
				"assume_role": map[string]any{"role_arn": "arn"},
				"features":    map[string]any{"manifest": true},
			},
			indent:  2,
			objects: []string{"assume_role"},
			want:    "  assume_role = {\n    role_arn = \"arn\"\n  }\n  features {\n    manifest = true\n  }\n  region = \"us-west-2\"\n",
		},
		{
			name: "objects at nested depth",
			config: map[string]any{
				"outer": map[string]any{
					"inner_obj": map[string]any{
						"key": "val",
					},
				},
			},
			indent:  2,
			objects: []string{"inner_obj"},
			want:    "  outer {\n    inner_obj = {\n      key = \"val\"\n    }\n  }\n",
		},
		{
			name: "sorted keys with blocks and objects",
			config: map[string]any{
				"region":      "us-west-2",
				"assume_role": map[string]any{"role_arn": "arn"},
			},
			indent:  2,
			objects: []string{"assume_role"},
			want:    "  assume_role = {\n    role_arn = \"arn\"\n  }\n  region = \"us-west-2\"\n",
		},
		{
			name:   "custom indent",
			config: map[string]any{"key": "val"},
			indent: 4,
			want:   "    key = \"val\"\n",
		},
		{
			name:    "nil objects list same as empty",
			config:  map[string]any{"block": map[string]any{"k": "v"}},
			indent:  2,
			objects: nil,
			want:    "  block {\n    k = \"v\"\n  }\n",
		},
		{
			name: "helm-like nested blocks",
			config: map[string]any{
				"kubernetes": map[string]any{
					"config_path": "~/.kube/config",
				},
				"experiments": map[string]any{
					"manifest": true,
				},
			},
			indent: 2,
			want:   "  experiments {\n    manifest = true\n  }\n  kubernetes {\n    config_path = \"~/.kube/config\"\n  }\n",
		},
		{
			name: "awscc-like attribute object with template-resolved value",
			config: map[string]any{
				"region": "us-east-1",
				"assume_role": map[string]any{
					"role_arn": "arn:aws:iam::123456789:role/tfe-si",
				},
			},
			indent:  2,
			objects: []string{"assume_role"},
			want:    "  assume_role = {\n    role_arn = \"arn:aws:iam::123456789:role/tfe-si\"\n  }\n  region = \"us-east-1\"\n",
		},
		{
			name: "azurerm-like empty features block",
			config: map[string]any{
				"resource_provider_registrations": "none",
				"features":                        map[string]any{},
			},
			indent: 2,
			want:   "  features {}\n  resource_provider_registrations = \"none\"\n",
		},
		{
			name: "dns-like triple nested blocks",
			config: map[string]any{
				"update": map[string]any{
					"server": "192.168.0.1",
					"port":   float64(53),
					"gssapi": map[string]any{
						"realm":    "EXAMPLE.COM",
						"username": "admin",
					},
				},
			},
			indent: 2,
			want:   "  update {\n    gssapi {\n      realm = \"EXAMPLE.COM\"\n      username = \"admin\"\n    }\n    port = 53\n    server = \"192.168.0.1\"\n  }\n",
		},
		{
			name: "multiple objects in list",
			config: map[string]any{
				"assume_role": map[string]any{"role_arn": "arn"},
				"tags":        map[string]any{"env": "prod"},
				"features":    map[string]any{"flag": true},
			},
			indent:  2,
			objects: []string{"assume_role", "tags"},
			want:    "  assume_role = {\n    role_arn = \"arn\"\n  }\n  features {\n    flag = true\n  }\n  tags = {\n    env = \"prod\"\n  }\n",
		},
		{
			name: "object propagates object syntax to nested maps",
			config: map[string]any{
				"outer_obj": map[string]any{
					"nested": map[string]any{
						"key": "val",
					},
					"scalar": "hello",
				},
			},
			indent:  2,
			objects: []string{"outer_obj"},
			want:    "  outer_obj = {\n    nested = {\n      key = \"val\"\n    }\n    scalar = \"hello\"\n  }\n",
		},
		{
			name: "block with attribute object inside",
			config: map[string]any{
				"outer_block": map[string]any{
					"nested_obj": map[string]any{
						"key": "val",
					},
					"scalar": "hello",
				},
			},
			indent:  2,
			objects: []string{"nested_obj"},
			want:    "  outer_block {\n    nested_obj = {\n      key = \"val\"\n    }\n    scalar = \"hello\"\n  }\n",
		},
		{
			name:   "float value",
			config: map[string]any{"ratio": float64(3.14)},
			indent: 2,
			want:   "  ratio = 3.14\n",
		},
		{
			name:   "single element list",
			config: map[string]any{"zones": []any{"us-east-1a"}},
			indent: 2,
			want:   "  zones = [\"us-east-1a\"]\n",
		},
		{
			name:   "empty list",
			config: map[string]any{"tags": []any{}},
			indent: 2,
			want:   "  tags = []\n",
		},
		{
			name: "google-like flat scalars",
			config: map[string]any{
				"project": "my-project-id",
				"region":  "us-central1",
			},
			indent: 2,
			want:   "  project = \"my-project-id\"\n  region = \"us-central1\"\n",
		},
		{
			name: "datadog-like block with nested attribute map",
			config: map[string]any{
				"api_key": "key123",
				"default_tags": map[string]any{
					"tags": map[string]any{
						"env":     "prod",
						"service": "web",
					},
				},
			},
			indent: 2,
			want:   "  api_key = \"key123\"\n  default_tags {\n    tags {\n      env = \"prod\"\n      service = \"web\"\n    }\n  }\n",
		},
		{
			name: "unrecognized objects key has no effect on other keys",
			config: map[string]any{
				"features": map[string]any{"flag": true},
			},
			indent:  2,
			objects: []string{"nonexistent"},
			want:    "  features {\n    flag = true\n  }\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderHCLBody(tt.config, tt.indent, tt.objects)
			require.Equal(t, tt.want, got)
		})
	}
}
