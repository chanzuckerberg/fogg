package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check for marker inconsistencies",
	Long:  `Scans for .grid-state.yaml files and checks for duplicates or conflicts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		markers, err := ScanMarkers(cwd)
		if err != nil {
			return err
		}

		issues := ValidateMarkers(markers)
		fmt.Printf("Found %d markers.\n", len(markers))
		if len(issues) > 0 {
			for _, issue := range issues {
				fmt.Printf("❌ %s\n", issue)
			}
			return fmt.Errorf("doctor found issues")
		}

		fmt.Println("✅ No issues found.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
