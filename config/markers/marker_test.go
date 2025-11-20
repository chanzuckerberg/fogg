package markers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chanzuckerberg/fogg/config/markers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadSaveMarker(t *testing.T) {
	tmpDir := t.TempDir()
	markerPath := filepath.Join(tmpDir, ".grid-state.yaml")

	originalMarker := &markers.Marker{
		GUID:      "test-guid",
		LogicalID: "test-logical-id",
		Labels: map[string]string{
			"env": "test",
		},
		Dependencies: []markers.Dependency{{GUID: "dep1"}, {GUID: "dep2", Output: "out", Input: "in"}},
	}

	err := markers.SaveMarker(markerPath, originalMarker)
	require.NoError(t, err)

	loadedMarker, err := markers.LoadMarker(markerPath)
	require.NoError(t, err)

	assert.Equal(t, originalMarker, loadedMarker)
}

func TestLoadMarker_NotFound(t *testing.T) {
	_, err := markers.LoadMarker("non-existent-file.yaml")
	assert.Error(t, err)
}

func TestSaveMarker_CreateDir(t *testing.T) {
	tmpDir := t.TempDir()
	markerPath := filepath.Join(tmpDir, "subdir", ".grid-state.yaml")

	marker := &markers.Marker{
		GUID: "test-guid",
	}

	err := markers.SaveMarker(markerPath, marker)
	require.NoError(t, err)

	_, err = os.Stat(markerPath)
	assert.NoError(t, err)
}
