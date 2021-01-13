package cmd

import (
	"github.com/chanzuckerberg/fogg/exp/migrate"
	"github.com/spf13/cobra"
)

func init() {
	expCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Assists with terraform state migrations",
	Long: `This command aims to assist with terraform state migrations.
	Particularly when there are module renames and such.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return migrate.Migrate("plan.tfplan")
	},
}
