package main

import (
	"fmt"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/gamestate"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/rs/zerolog/log"
)

func completedConstructions(userId string, gs *models.GameState, tx *gamestate.GameTx) error {
	completedConstructions := getCompletedConstructions(gs)

	// Exit early if there are no completed constructions
	if len(completedConstructions) == 0 {
		return nil
	}

	// Update patch
	u := tx.Users[userId]
	u.SetConstructionQueue(gs.ConstructionQueue[len(completedConstructions):])

	for _, constr := range completedConstructions {
		var lot *models.GameState_Lot
		var buildInfo *models.BuildingInfo
		var levelInfo *models.BuildingInfo_LevelInfo

		if constr.Razing {
			lot = gs.Lots[constr.LotId]
			u.RazeBuilding(constr.LotId)
		} else {
			u.ConstructBuilding(constr.LotId, constr.Building, constr.Level)
			lot = gs.Lots[constr.LotId]

			if u.LatestActivity.Before(time.Now().Add(-5 * time.Minute)) {
				u.AppendMessage(newConstructionCompletedMessage(userId, constr))
			}
		}

		buildInfo = internal.FullGameData.Buildings[int32(lot.Building)]
		if buildInfo == nil {
			log.Error().Int32("building", int32(lot.Building)).Msg("Could not find building info")
			continue
		}

		levelInfo = buildInfo.LevelInfos[lot.Level]
		if levelInfo == nil {
			log.Error().
				Int32("building", int32(lot.Building)).
				Int32("level", lot.Level).
				Msg("Could not find level info")
			continue
		}

		if levelInfo.Residence != nil {
			if !constr.Razing {
				var count int32
				if constr.Level > 0 {
					prevLevelInfo := buildInfo.LevelInfos[lot.Level-1]
					count = levelInfo.Residence.Beds - prevLevelInfo.Residence.Beds
				} else {
					count = levelInfo.Residence.Beds
				}
				for n := 0; n < int(count); n++ {
					u.AppendNewMouse()
				}
			} else {
				removeCount := levelInfo.Residence.Beds

				// Remove uneducated until there are not more uneducated mice
				for removeCount > 0 {
					if u.RemoveMouseByEducation(false, 0) {
						removeCount--
					} else {
						break
					}
				}

				// Remove educated mice until removeCount is 0
				popKey := 0
				loopCount := 0
				for removeCount > 0 && loopCount < 1000 {
					switch popKey {
					case 0:
						if u.RemoveMouseByEducation(true, models.Education_CHEF) {
							removeCount--
						}
					case 1:
						if u.RemoveMouseByEducation(true, models.Education_SALESMOUSE) {
							removeCount = removeCount - 1
						}
					case 2:
						if u.RemoveMouseByEducation(true, models.Education_GUARD) {
							removeCount = removeCount - 1
						}
					case 3:
						if u.RemoveMouseByEducation(true, models.Education_THIEF) {
							removeCount = removeCount - 1
						}
					case 4:
						if u.RemoveMouseByEducation(true, models.Education_PUBLICIST) {
							removeCount = removeCount - 1
						}
					}
					popKey = (popKey + 1) % 5
					loopCount++
				}
			}
		}
	}

	// Completion of buildings can affect the stats because we increase the
	// number of employables. E.g. if the player had 10 chefs but only 5 of them
	// were employed.
	u.StatsInvalidated = true

	return nil
}

func getCompletedConstructions(gs *models.GameState) (res []*models.Construction) {
	now := time.Now().UnixNano()

	for _, t := range gs.ConstructionQueue {
		if t.CompleteAt > now {
			break
		}

		res = append(res, t)
	}

	return res
}

func newConstructionCompletedMessage(userId string, constr *models.Construction) *messaging.Message {
	// internal.FullGameData.Buildings[int32(constr.Building)]
	var title string
	var body string

	buildInfo := internal.FullGameData.Buildings[int32(constr.Building)]
	if buildInfo != nil {
		title = "Construction completed!"
		if constr.Level > 0 {
			body = fmt.Sprintf("Your %s (level %d) has been completed.", buildInfo.Title, constr.Level)
		} else {
			body = fmt.Sprintf("Your %s has been completed.", buildInfo.Title)
		}
	} else {
		title = "Construction completed!"
	}

	return &messaging.Message{
		Data: map[string]string{
			"userId": userId,
		},
		Notification: &messaging.Notification{
			Title: title,
			Body: body,
		},
		Android: &messaging.AndroidConfig{
			CollapseKey: "reminder",
		},
		Webpush: &messaging.WebpushConfig{
			Notification: &messaging.WebpushNotification{
				Tag: "reminder",
			},
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-collapse-id": "reminder",
			},
		},
	}
}

