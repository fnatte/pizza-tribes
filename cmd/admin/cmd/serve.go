package cmd

import (
	"github.com/fnatte/pizza-tribes/cmd/admin/serve"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start admin server",
	Run: func(cmd *cobra.Command, args []string) {
		serve.Serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
