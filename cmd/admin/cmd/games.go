package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fnatte/pizza-tribes/cmd/admin/db"
	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/models"
	"github.com/fnatte/pizza-tribes/internal/mama"
	"github.com/spf13/cobra"
)

var title string
var host string

var gamesCmd = &cobra.Command{
	Use:   "games",
	Short: "Manage games",
}

var listGamesCmd = &cobra.Command{
	Use:   "ls",
	Short: "Get list of games",
	RunE: func(cmd *cobra.Command, args []string) error {
		sqldb := db.NewSqlClient()
		ctx := cmd.Context()

		games, err := mama.GetAllGames(ctx, sqldb)
		if err != nil {
			return err
		}

		fmt.Println("Games")
		fmt.Println("Id\tTitle\tStatus")
		for _, game := range games {
			fmt.Printf("%s\t%s\t%s\n", game.Id, game.Title, game.Status)
		}

		return nil
	},
}

var newGameCmd = &cobra.Command{
	Use:   "new",
	Short: "Start a new game",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		sqldb := db.NewSqlClient()
		gc := mama.NewGameCreator(sqldb)
		ctx := cmd.Context()
		if err := gc.CreateGame(ctx, title, host); err != nil {
			return err
		}

		fmt.Println("Game created")

		return nil
	},
}

var setGameStartTimeCmd = &cobra.Command{
	Use:   "starttime",
	Short: "Get or set the game start time",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rc := db.NewRedisClient()
		world := game.NewWorldService(rc)

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

			fmt.Printf("Start time set to %s\n", time.Unix(startTime, 0).Format(time.RFC3339))

			break
		}
	},
}



func init() {
	newGameCmd.PersistentFlags().StringVar(&title, "title", "", "game title")
	newGameCmd.MarkPersistentFlagRequired("title")
	newGameCmd.PersistentFlags().StringVar(&host, "host", "", "game host")
	newGameCmd.MarkPersistentFlagRequired("host")

	rootCmd.AddCommand(gamesCmd)

	gamesCmd.AddCommand(listGamesCmd)
	gamesCmd.AddCommand(newGameCmd)
	gamesCmd.AddCommand(setGameStartTimeCmd)
}
