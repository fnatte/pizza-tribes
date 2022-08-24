package internal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/rand"
)

func TestCalculateHeist(t *testing.T) {
	tests := map[string]struct {
		input Heist
		want  HeistOutcome
	}{
		"no guards, no loot": {
			input: Heist{
				Thieves:     1,
				TargetCoins: 0,
			},
			want: HeistOutcome{
				SuccessfulThieves: 1,
				Loot:              0,
			},
		},
		"no guards with all loot": {
			input: Heist{
				Thieves:     1,
				TargetCoins: 50,
			},
			want: HeistOutcome{
				SuccessfulThieves: 1,
				Loot:              50,
			},
		},
		"no guards with capacity loot": {
			input: Heist{
				Thieves:     1,
				TargetCoins: 100000,
			},
			want: HeistOutcome{
				SuccessfulThieves: 1,
				Loot:              3000,
			},
		},
		"10 thieves on 10 guards": {
			input: Heist{
				Thieves: 10,
				Guards:  10,
			},
			want: HeistOutcome{
				SuccessfulThieves: 7,
				CaughtThieves:     3,
				SleepingGuards:    1,
			},
		},
		"1 thief on 20 guards": {
			input: Heist{
				Thieves: 1,
				Guards:  20,
			},
			want: HeistOutcome{
				SuccessfulThieves: 0,
				CaughtThieves:     1,
				SleepingGuards:    1,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := CalculateHeist(test.input, rand.NewSource(0))

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CalculateHeist(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}

}
