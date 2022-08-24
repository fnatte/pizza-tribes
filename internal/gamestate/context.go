package gamestate

import (
	"fmt"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/rs/xid"
)

type GameTx_User struct {
	Gs             *models.GameState
	PatchMask      *models.PatchMask
	Reports        []*models.Report
	Messages       []*messaging.Message
	LatestActivity time.Time

	NextUpdateInvalidated bool
	StatsInvalidated      bool
	ReportsInvalidated    bool
	CoinsChanged          bool
	PizzasChanged         bool
}

type GameTx struct {
	Users map[string]*GameTx_User
}

func (u *GameTx_User) SetCoins(val int32) {
	if u.Gs.Resources.Coins != val {
		u.Gs.Resources.Coins = val
		u.CoinsChanged = true
		u.PatchMask.AppendPath("resources.coins")
	}
}

func (u *GameTx_User) SetPizzas(val int32) {
	if u.Gs.Resources.Pizzas != val {
		u.Gs.Resources.Pizzas = val
		u.PizzasChanged = true
		u.PatchMask.AppendPath("resources.pizzas")
	}
}

func (u *GameTx_User) IncrCoins(val int32) {
	u.SetCoins(u.Gs.Resources.Coins + val)
}

func (u *GameTx_User) IncrPizzas(val int32) {
	u.SetCoins(u.Gs.Resources.Pizzas + val)
}

func (u *GameTx_User) SetPizzaPrice(val int32) {
	if u.Gs.PizzaPrice != val {
		u.Gs.PizzaPrice = val
		u.StatsInvalidated = true
		u.PatchMask.AppendPath("pizzaPrice")
	}
}

func (u *GameTx_User) SetTimestamp(val int64) {
	if u.Gs.Timestamp != val {
		u.Gs.Timestamp = val
		u.PatchMask.AppendPath("timestamp")
	}
}

func (u *GameTx_User) SetTappedAt(lotId string, val int64) {
	if u.Gs.Lots[lotId].TappedAt != val {
		u.Gs.Lots[lotId].TappedAt = val
		u.PatchMask.AppendPath(fmt.Sprintf("lots.%s", lotId))
	}
}

func (u *GameTx_User) SetTaps(lotId string, val int32) {
	if u.Gs.Lots[lotId].Taps != val {
		u.Gs.Lots[lotId].Taps = val
		u.PatchMask.AppendPath(fmt.Sprintf("lots.%s", lotId))
	}
}

func (u *GameTx_User) SetStreak(lotId string, val int32) {
	if u.Gs.Lots[lotId].Streak != val {
		u.Gs.Lots[lotId].Streak = val
		u.PatchMask.AppendPath(fmt.Sprintf("lots.%s", lotId))
	}
}

func (u *GameTx_User) SetMouseIsBeingEducated(mouseId string, val bool) {
	if u.Gs.Mice[mouseId].IsBeingEducated != val {
		u.Gs.Mice[mouseId].IsBeingEducated = val
		u.PatchMask.AppendPath(fmt.Sprintf("mice.%s", mouseId))
	}
}

func (u *GameTx_User) SetMouseEducation(mouseId string, education models.Education) {
	m := u.Gs.Mice[mouseId]
	m.IsBeingEducated = false
	m.IsEducated = true
	m.Education = education
	u.PatchMask.AppendPath(fmt.Sprintf("mice.%s", mouseId))
}

func (u *GameTx_User) SetMouseUneducated(mouseId string) {
	m := u.Gs.Mice[mouseId]
	m.IsBeingEducated = false
	m.IsEducated = false
	m.Education = 0
	u.PatchMask.AppendPath(fmt.Sprintf("mice.%s", mouseId))
}

func (u *GameTx_User) SetMouseAppearance(mouseId string, appearance *models.MouseAppearance) {
	u.Gs.Mice[mouseId].Appearance = appearance
	u.PatchMask.AppendPath(fmt.Sprintf("mice.%s", mouseId))
}

func (u *GameTx_User) SetAmbassadorMouse(mouseId string) {
	u.Gs.AmbassadorMouseId = mouseId
	u.PatchMask.AppendPath("ambassadorMouseId")
}

func (u *GameTx_User) AppendNewMouse() {
	mouseId := xid.New().String()
	newMouse := &models.Mouse{
		Name:       internal.GetNewMouseName(u.Gs.Mice),
		IsEducated: false,
		Education:  0,
	}
	u.Gs.Mice[mouseId] = newMouse
	u.PatchMask.AppendPath(fmt.Sprintf("mice.%s", mouseId))
}

func (u *GameTx_User) RemoveMouseByEducation(isEducated bool, education models.Education) bool {
	for id, m := range u.Gs.Mice {
		if (!m.IsEducated && !isEducated) || m.Education == education {
			delete(u.Gs.Mice, id)
			u.PatchMask.AppendPath(fmt.Sprintf("mice.%s", id))
			return true
		}
	}

	return false
}

func (u *GameTx_User) SetTrainingQueue(val []*models.Training) {
	u.Gs.TrainingQueue = val
	u.PatchMask.AppendPath("trainingQueue")
}

func (u *GameTx_User) AppendTrainingQueue(val *models.Training) {
	u.Gs.TrainingQueue = append(u.Gs.TrainingQueue, val)
	u.NextUpdateInvalidated = true
	u.PatchMask.AppendPath("trainingQueue")
}

func (u *GameTx_User) SetConstructionQueue(val []*models.Construction) {
	u.Gs.ConstructionQueue = val
	u.PatchMask.AppendPath("constructionQueue")
}

func (u *GameTx_User) ConstructBuilding(lotId string, building models.Building, level int32) {
	lot := u.Gs.Lots[lotId]
	if lot != nil {
		u.Gs.Lots[lotId].Building = building
		u.Gs.Lots[lotId].Level = level
		u.PatchMask.AppendPath(fmt.Sprintf("lots.%s", lotId))
	} else {
		u.Gs.Lots[lotId] = &models.GameState_Lot{
			Building: building,
			Level:    level,
		}
		u.PatchMask.AppendPath(fmt.Sprintf("lots.%s", lotId))
	}
}

func (u *GameTx_User) RazeBuilding(lotId string) {
	lot := u.Gs.Lots[lotId]
	if lot != nil {
		delete(u.Gs.Lots, lotId)
		u.PatchMask.AppendPath(fmt.Sprintf("lots.%s", lotId))
	}
}

func (u *GameTx_User) SetTravelQueue(val []*models.Travel) {
	u.Gs.TravelQueue = val
	u.PatchMask.AppendPath("travelQueue")
}

func (u *GameTx_User) SetResearchQueue(val []*models.OngoingResearch) {
	u.Gs.ResearchQueue = val
	u.PatchMask.AppendPath("researchQueue")
}

func (u *GameTx_User) AppendDiscovery(val models.ResearchDiscovery) {
	u.Gs.Discoveries = append(u.Gs.Discoveries, val)
	u.PatchMask.AppendPath("discoveries")
}

func (u *GameTx_User) AppendReport(val *models.Report) {
	u.Reports = append(u.Reports, val)
	u.ReportsInvalidated = true
}

func (u *GameTx_User) AppendMessage(val *messaging.Message) {
	u.Messages = append(u.Messages, val)
}

func (u *GameTx_User) SetQuestCompleted(questId string) {
	if !u.Gs.Quests[questId].Completed {
		u.Gs.Quests[questId].Completed = true
		u.PatchMask.AppendPath(fmt.Sprintf("quests.%s.completed", questId))
	}
}

func (u *GameTx_User) SetQuestAvailable(questId string) {
	if u.Gs.Quests[questId] == nil {
		u.Gs.Quests[questId] = &models.QuestState{}
		u.PatchMask.AppendPath(fmt.Sprintf("quests.%s", questId))
	}
}

func (u *GameTx_User) SetGeniusFlashes(val int32) {
	if u.Gs.GeniusFlashes != val {
		u.Gs.GeniusFlashes = val
		u.PatchMask.AppendPath("geniusFlashes")
	}
}

func (u *GameTx_User) IncrGeniusFlashes(val int32) {
	u.SetGeniusFlashes(u.Gs.GeniusFlashes + val)
}

func (u *GameTx_User) ToServerMessage() *models.ServerMessage {
	return &models.ServerMessage{
		Id: xid.New().String(),
		Payload: &models.ServerMessage_StateChange{
			StateChange: &models.GameStatePatch{
				GameState: u.Gs,
				PatchMask: u.PatchMask,
			},
		},
	}
}

func NewGameTx(userId string, gs *models.GameState, latestActivity time.Time) *GameTx {
	return &GameTx{
		Users: map[string]*GameTx_User{
			userId: {
				LatestActivity: latestActivity,
				Reports:        []*models.Report{},
				Gs:             gs,
				PatchMask: &models.PatchMask{
					Paths: []string{},
				},
			},
		},
	}
}

func (tx *GameTx) InitUser(userId string, gs *models.GameState) {
	if tx.Users[userId] == nil {
		tx.Users[userId] = &GameTx_User{
			Reports: []*models.Report{},
			Gs:      gs,
			PatchMask: &models.PatchMask{
				Paths: []string{},
			},
		}
	}
}

