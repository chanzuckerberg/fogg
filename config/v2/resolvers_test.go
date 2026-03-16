package v2_test

import (
	"testing"

	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/stretchr/testify/require"
)

func TestResolveTfLint(t *testing.T) {
	r := require.New(t)
	tru := true
	fal := false

	data := []struct {
		def    *bool
		over   *bool
		output *bool
	}{
		{nil, nil, &fal},
		{nil, &tru, &tru},
		{nil, &fal, &fal},
		{&tru, nil, &tru},
		{&tru, &tru, &tru},
		{&tru, &fal, &fal},
		{&fal, nil, &fal},
		{&fal, &tru, &tru},
		{&fal, &fal, &fal},
	}
	for _, test := range data {
		tt := test
		t.Run("", func(t *testing.T) {
			def := v2.Common{Tools: &v2.Tools{TfLint: &v2.TfLint{Enabled: tt.def}}}
			over := v2.Common{Tools: &v2.Tools{TfLint: &v2.TfLint{Enabled: tt.over}}}
			result := v2.ResolveTfLint(def, over)
			r.Equal(tt.output, result.Enabled)
		})
	}
}

func TestResolveCustomProviders(t *testing.T) {
	r := require.New(t)

	t.Run("empty when no custom providers", func(t *testing.T) {
		c1 := v2.Common{}
		result := v2.ResolveCustomProviders(c1)
		r.Len(result, 0)
	})

	t.Run("single provider", func(t *testing.T) {
		c1 := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						CommonProvider: v2.CommonProvider{Version: util.StrPtr("~> 1.0")},
						Source:         util.StrPtr("hashicorp/awscc"),
					},
				},
			},
		}
		result := v2.ResolveCustomProviders(c1)
		r.Len(result, 1)
		r.Equal("hashicorp/awscc", *result["awscc"].Source)
		r.Equal("~> 1.0", *result["awscc"].Version)
	})

	t.Run("inheritance overrides version", func(t *testing.T) {
		defaults := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						CommonProvider: v2.CommonProvider{Version: util.StrPtr("~> 1.0")},
						Source:         util.StrPtr("hashicorp/awscc"),
					},
				},
			},
		}
		component := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						CommonProvider: v2.CommonProvider{Version: util.StrPtr("~> 2.0")},
					},
				},
			},
		}
		result := v2.ResolveCustomProviders(defaults, component)
		r.Len(result, 1)
		r.Equal("hashicorp/awscc", *result["awscc"].Source)
		r.Equal("~> 2.0", *result["awscc"].Version)
	})

	t.Run("multiple providers merge", func(t *testing.T) {
		defaults := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						CommonProvider: v2.CommonProvider{Version: util.StrPtr("~> 1.0")},
						Source:         util.StrPtr("hashicorp/awscc"),
					},
				},
			},
		}
		component := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"newrelic": {
						CommonProvider: v2.CommonProvider{Version: util.StrPtr("~> 3.0")},
						Source:         util.StrPtr("newrelic/newrelic"),
					},
				},
			},
		}
		result := v2.ResolveCustomProviders(defaults, component)
		r.Len(result, 2)
		r.Equal("hashicorp/awscc", *result["awscc"].Source)
		r.Equal("newrelic/newrelic", *result["newrelic"].Source)
	})

	t.Run("disabled provider excluded", func(t *testing.T) {
		fal := false
		c1 := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						CommonProvider: v2.CommonProvider{
							Version: util.StrPtr("~> 1.0"),
							Enabled: &fal,
						},
						Source: util.StrPtr("hashicorp/awscc"),
					},
				},
			},
		}
		result := v2.ResolveCustomProviders(c1)
		r.Len(result, 1)
		r.False(*result["awscc"].Enabled)
	})

	t.Run("config merges across commons", func(t *testing.T) {
		defaults := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						CommonProvider: v2.CommonProvider{Version: util.StrPtr("~> 1.0")},
						Source:         util.StrPtr("hashicorp/awscc"),
						Config: map[string]any{
							"region":  "us-east-1",
							"profile": "default",
						},
					},
				},
			},
		}
		component := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						Config: map[string]any{
							"region": "eu-west-1",
						},
					},
				},
			},
		}
		result := v2.ResolveCustomProviders(defaults, component)
		r.Len(result, 1)
		r.Equal("eu-west-1", result["awscc"].Config["region"])
		r.Equal("default", result["awscc"].Config["profile"])
	})

	t.Run("config with nested map", func(t *testing.T) {
		c1 := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						CommonProvider: v2.CommonProvider{Version: util.StrPtr("~> 1.0")},
						Source:         util.StrPtr("hashicorp/awscc"),
						Config: map[string]any{
							"assume_role": map[string]any{
								"role_arn": "arn:aws:iam::123:role/foo",
							},
						},
					},
				},
			},
		}
		result := v2.ResolveCustomProviders(c1)
		r.Len(result, 1)
		nested := result["awscc"].Config["assume_role"].(map[string]any)
		r.Equal("arn:aws:iam::123:role/foo", nested["role_arn"])
	})

	t.Run("objects carried through", func(t *testing.T) {
		c1 := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						CommonProvider: v2.CommonProvider{Version: util.StrPtr("~> 1.0")},
						Source:         util.StrPtr("hashicorp/awscc"),
						Objects:        []string{"assume_role"},
					},
				},
			},
		}
		result := v2.ResolveCustomProviders(c1)
		r.Equal([]string{"assume_role"}, result["awscc"].Objects)
	})

	t.Run("objects replaced by later commons", func(t *testing.T) {
		defaults := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						CommonProvider: v2.CommonProvider{Version: util.StrPtr("~> 1.0")},
						Source:         util.StrPtr("hashicorp/awscc"),
						Objects:        []string{"assume_role"},
					},
				},
			},
		}
		component := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						Objects: []string{"assume_role", "extra"},
					},
				},
			},
		}
		result := v2.ResolveCustomProviders(defaults, component)
		r.Equal([]string{"assume_role", "extra"}, result["awscc"].Objects)
	})

	t.Run("objects nil does not clear existing", func(t *testing.T) {
		defaults := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						CommonProvider: v2.CommonProvider{Version: util.StrPtr("~> 1.0")},
						Source:         util.StrPtr("hashicorp/awscc"),
						Objects:        []string{"assume_role"},
					},
				},
			},
		}
		component := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {},
				},
			},
		}
		result := v2.ResolveCustomProviders(defaults, component)
		r.Equal([]string{"assume_role"}, result["awscc"].Objects)
	})

	t.Run("many providers from one common", func(t *testing.T) {
		c1 := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						CommonProvider: v2.CommonProvider{Version: util.StrPtr("~> 1.0")},
						Source:         util.StrPtr("hashicorp/awscc"),
						Objects:        []string{"assume_role"},
					},
					"azurerm": {
						CommonProvider: v2.CommonProvider{Version: util.StrPtr("~> 4.0")},
						Source:         util.StrPtr("hashicorp/azurerm"),
						Config:         map[string]any{"features": map[string]any{}},
					},
					"google": {
						CommonProvider: v2.CommonProvider{Version: util.StrPtr("~> 6.0")},
						Source:         util.StrPtr("hashicorp/google"),
						Config:         map[string]any{"project": "my-proj"},
					},
				},
			},
		}
		result := v2.ResolveCustomProviders(c1)
		r.Len(result, 3)
		r.Equal("hashicorp/awscc", *result["awscc"].Source)
		r.Equal([]string{"assume_role"}, result["awscc"].Objects)
		r.Equal("hashicorp/azurerm", *result["azurerm"].Source)
		r.Nil(result["azurerm"].Objects)
		features := result["azurerm"].Config["features"].(map[string]any)
		r.Len(features, 0)
		r.Equal("hashicorp/google", *result["google"].Source)
		r.Equal("my-proj", result["google"].Config["project"])
	})

	t.Run("config override does not mutate original", func(t *testing.T) {
		defaults := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						CommonProvider: v2.CommonProvider{Version: util.StrPtr("~> 1.0")},
						Source:         util.StrPtr("hashicorp/awscc"),
						Config:         map[string]any{"region": "us-east-1"},
					},
				},
			},
		}
		originalRegion := defaults.Providers.Custom["awscc"].Config["region"]
		component := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						Config: map[string]any{"region": "eu-west-1"},
					},
				},
			},
		}
		result := v2.ResolveCustomProviders(defaults, component)
		r.Equal("eu-west-1", result["awscc"].Config["region"])
		r.Equal("us-east-1", originalRegion)
	})

	t.Run("three level inheritance", func(t *testing.T) {
		defaults := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						CommonProvider: v2.CommonProvider{Version: util.StrPtr("~> 1.0")},
						Source:         util.StrPtr("hashicorp/awscc"),
						Config:         map[string]any{"region": "us-east-1", "profile": "default"},
						Objects:        []string{"assume_role"},
					},
				},
			},
		}
		env := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						CommonProvider: v2.CommonProvider{Version: util.StrPtr("~> 1.5")},
						Config:         map[string]any{"profile": "staging"},
					},
				},
			},
		}
		component := v2.Common{
			Providers: &v2.Providers{
				Custom: map[string]*v2.CustomProvider{
					"awscc": {
						Config: map[string]any{"region": "eu-west-1"},
					},
				},
			},
		}
		result := v2.ResolveCustomProviders(defaults, env, component)
		r.Len(result, 1)
		r.Equal("hashicorp/awscc", *result["awscc"].Source)
		r.Equal("~> 1.5", *result["awscc"].Version)
		r.Equal("eu-west-1", result["awscc"].Config["region"])
		r.Equal("staging", result["awscc"].Config["profile"])
		r.Equal([]string{"assume_role"}, result["awscc"].Objects)
	})
}
