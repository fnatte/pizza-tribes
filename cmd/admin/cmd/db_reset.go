package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/fnatte/pizza-tribes/cmd/admin/db"
)

func envOrDefault(key string, defaultVal string) string{
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return defaultVal
}

var dbResetCmd = &cobra.Command{
	Use: "reset",
	Short: "Reset state of server",
	Long: "Reset state of server. This will delete game state and users.",
	Run: func(cmd *cobra.Command, args []string) {
		rdb := db.NewRedisClient()

		res := rdb.FlushDB(cmd.Context())

		if res.Err() != nil {
			fmt.Println("An error occurred:")
			fmt.Println(res.Err().Error())
			os.Exit(1)
			return
		}

		fmt.Printf("Result: %s\n", res.String())
	},
}

func init() {
	dbCmd.AddCommand(dbResetCmd)
}

