package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fnatte/pizza-tribes/cmd/admin/db"
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/spf13/cobra"
)

var roundStartTimeCmd = &cobra.Command{
	Use:   "starttime",
	Short: "Get or set the game round start time",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rc := db.NewRedisClient()
		world := internal.NewWorldService(rc)

		state, err := world.GetState(cmd.Context())
		if err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
			return
		}

		// Get start time
		if len(args) == 0 {
			switch state.Type.(type) {
			case *models.WorldState_Started_:
				fmt.Printf("Game round has started\n")
				break
			case *models.WorldState_Starting_:
				fmt.Printf("Game round is starting at %d\n", state.StartTime)
				break
			case *models.WorldState_Ended_:
				fmt.Printf("Game round has ended\n")
				break
			}
			return
		}

		// Set start time
		switch state.Type.(type) {
		case *models.WorldState_Started_:
			fmt.Printf("Can not set start time because the round is already started\n")
			break
		case *models.WorldState_Ended_:
			fmt.Printf("Can not set start time because the round has already ended\n")
			break
		case *models.WorldState_Starting_:
			startTime, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				dur, err := time.ParseDuration(args[0])
				if err != nil {
					fmt.Printf("Could not parse start time")
					os.Exit(1)
					return
				}

				startTime = time.Now().Add(dur).Unix()
			}

			err = world.SetStartTime(cmd.Context(), startTime)
			if err != nil {
				fmt.Printf("Error: %s", err)
				os.Exit(1)
				return
			}
			break
		}
	},
}

func init() {
	roundCmd.AddCommand(roundStartTimeCmd)
}
