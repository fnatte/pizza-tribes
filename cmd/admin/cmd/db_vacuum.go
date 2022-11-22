package cmd

import (
	"github.com/fnatte/pizza-tribes/cmd/admin/db"
	"github.com/fnatte/pizza-tribes/internal/game/persist"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var dbVacuumCmd = &cobra.Command{
	Use: "vacuum",
	Short: "Vacuum state of server",
	Long: "Vacuum state of server. This will hanging game state.",
	RunE: func(cmd *cobra.Command, args []string) error {
		rc := db.NewRedisClient()
		ctx := cmd.Context()

		userRepo := persist.NewGameUserRepository(rc)

		userIds, err := userRepo.GetAllUsers(ctx)
		if err != nil {
			return err
		}

		res, err := rc.ZRevRange(ctx, "demands", 0, -1).Result()
		if err != nil {
			return err
		}
		for _, userId := range res {
			exists := false
			for _, uid := range userIds {
				if uid == userId {
					exists = true
					break
				}
			}

			if !exists {
				log.Info().Msgf("Delete user %s", userId)
				if err := userRepo.DeleteUser(cmd.Context(), userId); err != nil {
					return err
				}
			}
		}

		return nil
	},
}

func init() {
	dbCmd.AddCommand(dbVacuumCmd)
}

