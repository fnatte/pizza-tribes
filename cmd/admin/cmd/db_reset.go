package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/fnatte/pizza-tribes/cmd/admin/db"
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/spf13/cobra"
)

func envOrDefault(key string, defaultVal string) string{
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return defaultVal
}

func ensureWorld(ctx context.Context, r internal.RedisClient) error {
	world := internal.NewWorldService(r)
	if err := world.Initialize(ctx); err != nil {
		return err
	}

	return nil
}

var dbResetCmd = &cobra.Command{
	Use: "reset",
	Short: "Reset state of server",
	Long: "Reset state of server. This will delete game state and users.",
	Run: func(cmd *cobra.Command, args []string) {
		rdb := db.NewRedisClient()
		rc := internal.NewRedisClient(rdb)

		res := rdb.FlushDB(cmd.Context())
		if res.Err() != nil {
			fmt.Println("An error occurred:")
			fmt.Println(res.Err().Error())
			os.Exit(1)
			return
		}

		fmt.Printf("Flush result: %s\n", res.String())

		err := ensureWorld(cmd.Context(), rc)
		if err != nil {
			fmt.Printf("Failed to ensure world\n")
		}
		fmt.Printf("Initialized new world\n")

	},
}

func init() {
	dbCmd.AddCommand(dbResetCmd)
}

