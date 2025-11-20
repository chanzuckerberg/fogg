package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/chanzuckerberg/fogg/config/markers"
)

// LoadedMarker contains the marker data and its source path.
type LoadedMarker struct {
	Path   string
	Marker *markers.Marker
}

// ScanMarkers finds all .grid-state.yaml files in the current directory and subdirectories.
func ScanMarkers(root string) ([]LoadedMarker, error) {
	fsys := os.DirFS(root)
	// Find all .grid-state.yaml files
	matches, err := doublestar.Glob(fsys, "**/.grid-state.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to glob markers: %w", err)
	}

	var results []LoadedMarker
	for _, match := range matches {
		fullPath := filepath.Join(root, match)
		marker, err := markers.LoadMarker(fullPath)
		if err != nil {
			// We might want to warn instead of failing hard, but for now fail.
			return nil, fmt.Errorf("failed to load marker at %s: %w", fullPath, err)
		}
		results = append(results, LoadedMarker{
			Path:   fullPath,
			Marker: marker,
		})
	}

	return results, nil
}

// ValidateMarkers checks for required fields and conflicts.
// Returns a slice of human-readable issues. Empty slice means no issues found.
func ValidateMarkers(markers []LoadedMarker) []string {
	var issues []string
	guidMap := make(map[string][]LoadedMarker)
	logicIDMap := make(map[string][]LoadedMarker)

	for _, m := range markers {
		if strings.TrimSpace(m.Marker.GUID) == "" {
			issues = append(issues, fmt.Sprintf("%s is missing guid", m.Path))
		}
		if strings.TrimSpace(m.Marker.LogicalID) == "" {
			issues = append(issues, fmt.Sprintf("%s is missing logicalId", m.Path))
		}

		guidMap[m.Marker.GUID] = append(guidMap[m.Marker.GUID], m)
		logicIDMap[m.Marker.LogicalID] = append(logicIDMap[m.Marker.LogicalID], m)

		for i, dep := range m.Marker.Dependencies {
			if strings.TrimSpace(dep.GUID) == "" {
				issues = append(issues, fmt.Sprintf("%s dependency[%d] is missing guid", m.Path, i))
			}
			// Warn when output is not specified - will default to "default" in grid-sync
			if strings.TrimSpace(dep.Output) == "" {
				issues = append(issues, fmt.Sprintf("⚠️  %s dependency[%d] (guid=%s) has no output specified - will use 'default' output",
					m.Path, i, dep.GUID))
			}
		}
	}

	for guid, entries := range guidMap {
		if guid == "" {
			continue
		}
		if len(entries) > 1 {
			issues = append(issues, fmt.Sprintf("duplicate guid '%s' in %s", guid, joinPaths(entries)))
		}
	}

	for logicID, entries := range logicIDMap {
		if logicID == "" {
			continue
		}
		if len(entries) > 1 {
			issues = append(issues, fmt.Sprintf("duplicate logicalId '%s' in %s", logicID, joinPaths(entries)))
		}
	}

	return issues
}

func joinPaths(markers []LoadedMarker) string {
	paths := make([]string, 0, len(markers))
	for _, m := range markers {
		paths = append(paths, m.Path)
	}
	return strings.Join(paths, ", ")
}
