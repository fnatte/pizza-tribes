package cmd

import (
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Run commands against the db",
}

func init() {
	rootCmd.AddCommand(dbCmd)
}
