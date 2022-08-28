package game

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestReminderNextOccurence(t *testing.T) {
	tests := map[string]struct {
		reminder Reminder
		want     time.Time
		now      time.Time
	}{
		"simple": {
			reminder: Reminder{
				Interval: time.Hour,
				Offset:   time.Minute * 15,
			},
			now:  time.Date(2020, time.February, 10, 15, 20, 0, 0, time.UTC),
			want: time.Date(2020, time.February, 10, 16, 15, 0, 0, time.UTC),
		},
		"very soon": {
			reminder: Reminder{
				Interval: time.Hour,
				Offset:   time.Minute * 15,
			},
			now:  time.Date(2020, time.February, 10, 15, 14, 0, 0, time.UTC),
			want: time.Date(2020, time.February, 10, 15, 15, 0, 0, time.UTC),
		},
		"next day": {
			reminder: Reminder{
				Interval: 2 * time.Hour,
				Offset:   time.Minute * 30,
			},
			now:  time.Date(2020, time.February, 10, 22, 40, 0, 0, time.UTC),
			want: time.Date(2020, time.February, 11, 00, 30, 0, 0, time.UTC),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := test.reminder.NextOccurenceAfter(test.now)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("NextOccurenceAfter(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}

}
