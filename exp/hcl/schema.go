package hcl

import hcl "github.com/hashicorp/hcl/v2"

// taken from https://github.com/hashicorp/terraform-config-inspect/blob/17f92b0546e8602e4cecb14d2263f0b1746b9cc9/tfconfig/schema.go

var rootSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "terraform",
			LabelNames: nil,
		},
		{
			Type:       "variable",
			LabelNames: []string{"name"},
		},
		{
			Type:       "output",
			LabelNames: []string{"name"},
		},
		{
			Type:       "provider",
			LabelNames: []string{"name"},
		},
		{
			Type:       "resource",
			LabelNames: []string{"type", "name"},
		},
		{
			Type:       "data",
			LabelNames: []string{"type", "name"},
		},
		{
			Type:       "module",
			LabelNames: []string{"name"},
		},
		{
			Type: "locals",
		},
	},
}
