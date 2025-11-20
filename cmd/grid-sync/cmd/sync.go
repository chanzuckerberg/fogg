package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"connectrpc.com/connect"
	"github.com/chanzuckerberg/fogg/config/markers"
	"github.com/spf13/cobra"
	"github.com/terraconstructs/grid/pkg/sdk"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize markers to Grid API",
	Long:  `Reads .grid-state.yaml files and updates the Grid API state.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		markers, err := ScanMarkers(cwd)
		if err != nil {
			return err
		}
		if len(markers) == 0 {
			return fmt.Errorf("no markers found under %s", cwd)
		}

		if issues := ValidateMarkers(markers); len(issues) > 0 {
			for _, issue := range issues {
				fmt.Printf("❌ %s\n", issue)
			}
			return fmt.Errorf("marker validation failed")
		}

		session := sessionConfig{
			ServerURL:    opts.serverURL,
			ClientID:     opts.clientID,
			ClientSecret: opts.clientSecret,
		}

		client, err := newGridClient(ctx, session)
		if err != nil {
			return err
		}

		fmt.Printf("Synchronizing %d markers to %s...\n", len(markers), session.ServerURL)
		for _, marker := range markers {
			if err := syncMarker(ctx, client, marker); err != nil {
				return err
			}
		}

		fmt.Println("✅ Sync complete")
		return nil
	},
}

func syncMarker(ctx context.Context, client *sdk.Client, marker LoadedMarker) error {
	desiredLabels := markerLabelsToLabelMap(marker.Marker.Labels)
	desiredDeps := marker.Marker.Dependencies

	info, err := client.GetStateInfo(ctx, sdk.StateReference{GUID: marker.Marker.GUID})
	if err != nil {
		var cerr *connect.Error
		if errors.As(err, &cerr) && cerr.Code() == connect.CodeNotFound {
			fmt.Printf("[create] guid=%s logicalId=%s (%s)\n", marker.Marker.GUID, marker.Marker.LogicalID, marker.Path)
			_, createErr := client.CreateState(ctx, sdk.CreateStateInput{
				GUID:    marker.Marker.GUID,
				LogicID: marker.Marker.LogicalID,
				Labels:  desiredLabels,
			})
			if createErr != nil {
				return fmt.Errorf("failed to create state guid=%s: %w", marker.Marker.GUID, createErr)
			}
			return nil
		}
		return fmt.Errorf("failed to fetch state guid=%s: %w", marker.Marker.GUID, err)
	}

	if info.State.LogicID != marker.Marker.LogicalID {
		return fmt.Errorf("state guid=%s logicalId mismatch (grid=%s marker=%s): rename not supported", marker.Marker.GUID, info.State.LogicID, marker.Marker.LogicalID)
	}

	adds, removals := diffLabels(info.Labels, desiredLabels)
	if len(adds) == 0 && len(removals) == 0 {
		fmt.Printf("[noop-labels] guid=%s logicalId=%s\n", marker.Marker.GUID, marker.Marker.LogicalID)
	}

	_, err = client.UpdateStateLabels(ctx, sdk.UpdateStateLabelsInput{
		StateID:  info.State.GUID,
		Adds:     adds,
		Removals: removals,
	})
	if err != nil {
		return fmt.Errorf("failed to update labels for guid=%s: %w", info.State.GUID, err)
	}

	fmt.Printf("[labels] guid=%s adds=%d removals=%d\n", info.State.GUID, len(adds), len(removals))

	if err := syncDependencies(ctx, client, marker, info, desiredDeps); err != nil {
		return err
	}

	return nil
}

func markerLabelsToLabelMap(labels map[string]string) sdk.LabelMap {
	if len(labels) == 0 {
		return nil
	}
	result := make(sdk.LabelMap, len(labels))
	for k, v := range labels {
		result[k] = v
	}
	return result
}

func diffLabels(current sdk.LabelMap, desired sdk.LabelMap) (sdk.LabelMap, []string) {
	adds := make(sdk.LabelMap)
	var removals []string

	currentClean := make(map[string]any)
	for k, v := range current {
		currentClean[k] = v
	}

	for k, v := range desired {
		if cur, ok := currentClean[k]; !ok || !valuesEqual(cur, v) {
			adds[k] = v
		}
	}

	for k := range currentClean {
		if desired == nil {
			removals = append(removals, k)
			continue
		}
		if _, ok := desired[k]; !ok {
			removals = append(removals, k)
		}
	}

	return adds, removals
}

func valuesEqual(a any, b any) bool {
	switch va := a.(type) {
	case string:
		vb, ok := b.(string)
		return ok && va == vb
	case fmt.Stringer:
		vb, ok := b.(fmt.Stringer)
		return ok && va.String() == vb.String()
	default:
		return fmt.Sprint(a) == fmt.Sprint(b)
	}
}

func syncDependencies(ctx context.Context, client *sdk.Client, marker LoadedMarker, info *sdk.StateInfo, desired []markers.Dependency) error {
	const defaultOutput = "default"

	desiredMap := make(map[string]markers.Dependency)
	for _, dep := range desired {
		guid := strings.TrimSpace(dep.GUID)
		if guid == "" {
			continue
		}
		output := dep.Output
		if output == "" {
			output = defaultOutput
		}
		key := depKey(guid, output, dep.Input)
		desiredMap[key] = dep
	}

	existing := make(map[string]sdk.DependencyEdge)
	for _, edge := range info.Dependencies {
		// Only consider incoming edges for this state
		if edge.To.GUID != "" && edge.To.GUID != info.State.GUID {
			continue
		}
		if edge.To.GUID == "" && info.State.LogicID != "" && edge.To.LogicID != info.State.LogicID {
			continue
		}

		fromGUID := edge.From.GUID
		if fromGUID == "" {
			fromGUID = edge.From.LogicID
		}
		if fromGUID == "" {
			continue
		}
		output := edge.FromOutput
		if output == "" {
			output = defaultOutput
		}
		key := depKey(fromGUID, output, edge.ToInputName)
		existing[key] = edge
	}

	// Remove stale edges
	for key, edge := range existing {
		if _, ok := desiredMap[key]; ok {
			continue
		}
		if err := client.RemoveDependency(ctx, edge.ID); err != nil {
			return fmt.Errorf("failed to remove dependency %d: %w", edge.ID, err)
		}
		fmt.Printf("[deps-remove] guid=%s removed edge from=%s output=%s input=%s\n", info.State.GUID, edge.From.GUID, edge.FromOutput, edge.ToInputName)
	}

	// Add missing edges
	for key, dep := range desiredMap {
		if _, ok := existing[key]; ok {
			continue
		}
		output := dep.Output
		if output == "" {
			output = defaultOutput
		}
		input := dep.Input
		_, err := client.AddDependency(ctx, sdk.AddDependencyInput{
			From:        sdk.StateReference{GUID: dep.GUID},
			FromOutput:  output,
			To:          sdk.StateReference{GUID: info.State.GUID},
			ToInputName: input,
		})
		if err != nil {
			return fmt.Errorf("failed to add dependency from %s to %s: %w", dep.GUID, info.State.GUID, err)
		}
		fmt.Printf("[deps-add] guid=%s from=%s output=%s input=%s\n", info.State.GUID, dep.GUID, output, input)
	}

	return nil
}

func depKey(fromGUID, output, input string) string {
	return fromGUID + "|" + output + "|" + input
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
