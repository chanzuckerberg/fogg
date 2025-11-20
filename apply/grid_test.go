package apply

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/chanzuckerberg/fogg/config/markers"
	v2 "github.com/chanzuckerberg/fogg/config/v2"
	"github.com/chanzuckerberg/fogg/plan"
	"github.com/chanzuckerberg/fogg/util"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestResolveGridGUIDs_PreservesExisting(t *testing.T) {
	fs := afero.NewMemMapFs()
	enabled := true

	// Setup existing marker
	accountName := "foo"
	existingGUID := "existing-guid-123"
	markerPath := filepath.Join(util.RootPath, "accounts", accountName, ".grid-state.yaml")

	marker := markers.Marker{
		GUID:      existingGUID,
		LogicalID: "foo",
	}
	data, err := yaml.Marshal(marker)
	require.NoError(t, err)
	require.NoError(t, afero.WriteFile(fs, markerPath, data, 0644))

	// Setup Plan
	p := &plan.Plan{
		Accounts: map[string]plan.Account{
			accountName: {
				ComponentCommon: plan.ComponentCommon{
					Grid: &v2.GridConfig{
						Enabled: &enabled,
					},
				},
			},
		},
	}

	// Execute
	guids, err := resolveGridGUIDs(fs, p)
	require.NoError(t, err)

	// Verify
	assert.Equal(t, existingGUID, guids[fmt.Sprintf("account:%s", accountName)])
}

func TestResolveGridGUIDs_GeneratesNew(t *testing.T) {
	fs := afero.NewMemMapFs()
	enabled := true
	accountName := "bar"

	// Setup Plan (no existing marker)
	p := &plan.Plan{
		Accounts: map[string]plan.Account{
			accountName: {
				ComponentCommon: plan.ComponentCommon{
					Grid: &v2.GridConfig{
						Enabled: &enabled,
					},
				},
			},
		},
	}

	// Execute
	guids, err := resolveGridGUIDs(fs, p)
	require.NoError(t, err)

	// Verify
	guid := guids[fmt.Sprintf("account:%s", accountName)]
	assert.NotEmpty(t, guid)
	assert.NotEqual(t, "existing-guid-123", guid) // Just to be sure
}

func TestResolveGridGUIDs_OverrideTakesPrecedence(t *testing.T) {
	fs := afero.NewMemMapFs()
	enabled := true
	overrideGUID := "override-guid-456"
	existingGUID := "existing-guid-123"
	accountName := "baz"

	// Setup existing marker
	markerPath := filepath.Join(util.RootPath, "accounts", accountName, ".grid-state.yaml")
	marker := markers.Marker{
		GUID:      existingGUID,
		LogicalID: "baz",
	}
	data, err := yaml.Marshal(marker)
	require.NoError(t, err)
	require.NoError(t, afero.WriteFile(fs, markerPath, data, 0644))

	// Setup Plan with Override
	p := &plan.Plan{
		Accounts: map[string]plan.Account{
			accountName: {
				ComponentCommon: plan.ComponentCommon{
					Grid: &v2.GridConfig{
						Enabled: &enabled,
						GUID:    &overrideGUID,
					},
				},
			},
		},
	}

	// Execute
	guids, err := resolveGridGUIDs(fs, p)
	require.NoError(t, err)

	// Verify
	assert.Equal(t, overrideGUID, guids[fmt.Sprintf("account:%s", accountName)])
}
