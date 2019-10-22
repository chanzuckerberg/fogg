package exp

import (
	"fmt"
	"os"

	"github.com/go-logfmt/logfmt"
	"github.com/hashicorp/terraform/plans"
	"github.com/hashicorp/terraform/plans/planfile"
	"github.com/spf13/cobra"
)

func init() {
	entropyCmd.Flags().StringP("plan-file", "f", "TODO", "Path to Terraform Plan file to parse.")
	entropyCmd.Flags().StringP("output-file", "o", "TODO", "Path to write instrumentation to")

	ExpCmd.AddCommand(entropyCmd)
}

const (
	allActions = "AllActions"
)

type terraformDiff struct {
	Address      string `logfmt:"address,omitempty"`
	ResourceMode string `logfmt:"resource_mode,omitempty"`
	Action       string `logfmt:"action,omitempty"`
	Project      string `logfmt:"project,omitempty"`
	Component    string `logfmt:"component,omitempty"`
}

var entropyCmd = &cobra.Command{
	Use:   "entropy",
	Short: "Measures how many differences result from a terraform plan.",
	Long: `This command will parse a Terraform plan and track any diffs.
It is meant to be run with honeycomb/buildevents and thus we generate
output in LogFmt format.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := entropyRun(cmd, args)
		if err != nil {
			// We don't want this to error out our build
			fmt.Fprintf(os.Stderr, "fogg entropy error: %s", err)
		}
	},
}

func entropyRun(cmd *cobra.Command, args []string) error {
	planFilePath, err := cmd.Flags().GetString("plan-file")
	if err != nil {
		return fmt.Errorf("could not read plan-file flag: %w", err)
	}
	outputFilePath, err := cmd.Flags().GetString("output-file")
	if err != nil {
		return fmt.Errorf("could not read output-file flag: %w", err)
	}

	planReader, err := planfile.Open(planFilePath)
	if err != nil {
		return fmt.Errorf("could not open terraform plan: %w", err)
	}
	defer planReader.Close()

	plan, err := planReader.ReadPlan()
	if err != nil {
		return fmt.Errorf("could not read/parse terraform plan: %w", err)
	}

	f, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("could not open output file for writing: %w", err)
	}
	defer f.Close()

	encoder := logfmt.NewEncoder(f)
	actionCounts := map[string]int{
		// Also keep a summary of all actions
		allActions: 0,
	}

	// We just keep a simple count of the terraform actions
	for _, resourceChange := range plan.Changes.Resources {
		action := resourceChange.Action

		if action == plans.NoOp {
			continue
		}

		count, ok := actionCounts[action.String()]
		if !ok {
			actionCounts[action.String()] = 1
		} else {
			actionCounts[action.String()] = count + 1
		}
		actionCounts[allActions] = actionCounts[allActions] + 1
	}

	for action, count := range actionCounts {
		err = encoder.EncodeKeyval(action, count)
		if err != nil {
			return fmt.Errorf("could not encode key/val %w", err)
		}
	}

	err = encoder.EndRecord()
	if err != nil {
		return fmt.Errorf("could not end record: %w", err)
	}
	return nil
}
