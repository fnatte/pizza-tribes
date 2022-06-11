package cmd

import (
	"github.com/spf13/cobra"
)

// roundCmd represents the round command
var roundCmd = &cobra.Command{
	Use:   "round",
	Short: "Manage game round",
}

func init() {
	rootCmd.AddCommand(roundCmd)
}
