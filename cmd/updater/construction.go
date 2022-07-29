package main

import (
	"fmt"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/gamestate"
	"github.com/fnatte/pizza-tribes/internal/models"
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
		if constr.Razing {
			u.RazeBuilding(constr.LotId)
		} else {
			u.ConstructBuilding(constr.LotId, constr.Building, constr.Level)

			if u.LatestActivity.Before(time.Now().Add(-5 * time.Minute)) {
				u.AppendMessage(newConstructionCompletedMessage(userId, constr))
			}
		}

		buildInfo := internal.FullGameData.Buildings[int32(constr.Building)]
		if buildInfo == nil {
			continue
		}

		levelInfo := buildInfo.LevelInfos[constr.Level]
		if levelInfo == nil {
			continue
		}

		if levelInfo.Residence != nil {
			if !constr.Razing {
				var count int32
				if constr.Level > 0 {
					prevLevelInfo := buildInfo.LevelInfos[constr.Level-1]
					count = levelInfo.Residence.Beds - prevLevelInfo.Residence.Beds
				} else {
					count = levelInfo.Residence.Beds
				}
				u.IncrUneducated(count)
				for n := 0; n < int(count); n++ {
					u.AppendNewMouse()
				}
			} else {
				rest := levelInfo.Residence.Beds

				if rest > gs.Population.Uneducated {
					u.SetUneducated(0)
					rest = rest - gs.Population.Uneducated
					for n := 0; n < int(gs.Population.Uneducated); n++ {
						u.RemoveMouseByEducation(false, 0)
					}
				} else {
					u.IncrUneducated(-rest)
					rest = 0
					for n := 0; n < int(rest); n++ {
						u.RemoveMouseByEducation(false, 0)
					}
				}

				popKey := 0
				loopCount := 0
				for rest > 0 && loopCount < 1000 {
					switch popKey {
					case 0:
						if gs.Population.Chefs > 0 {
							u.IncrChefs(-1)
							u.RemoveMouseByEducation(true, models.Education_CHEF)
							rest = rest - 1
						}
					case 1:
						if gs.Population.Salesmice > 0 {
							u.IncrSalesmice(-1)
							u.RemoveMouseByEducation(true, models.Education_SALESMOUSE)
							rest = rest - 1
						}
					case 2:
						if gs.Population.Guards > 0 {
							u.IncrGuards(-1)
							u.RemoveMouseByEducation(true, models.Education_GUARD)
							rest = rest - 1
						}
					case 3:
						if gs.Population.Thieves > 0 {
							u.IncrThieves(-1)
							u.RemoveMouseByEducation(true, models.Education_THIEF)
							rest = rest - 1
						}
					case 4:
						if gs.Population.Publicists > 0 {
							u.IncrPublicists(-1)
							u.RemoveMouseByEducation(true, models.Education_PUBLICIST)
							rest = rest - 1
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

