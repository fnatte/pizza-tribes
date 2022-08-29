package game

import (
	"testing"
)

func TestQuestOrderIsUnique(t *testing.T) {
	for i, q1 := range FullGameData.Quests {
		for j, q2 := range FullGameData.Quests {
			if i != j && q1.Order == q2.Order {
				t.Errorf("Quest %d and %d has the same order %d", i, j, q1.Order)
			}
		}
	}
}

func TestEveryQuestOneOfItemsExist(t *testing.T) {
	for _, q := range FullGameData.Quests {
		for _, id := range q.Reward.OneOfItems {
			if _, ok := AppearancePartsMap[id]; !ok {
				t.Errorf("Quest %s has an item reward that does not exist: %s", q.Id, id)
			}
		}
	}
}
