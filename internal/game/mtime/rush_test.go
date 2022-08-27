package mtime

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type input struct{ then int64; now int64 }
type result struct{ Rush int64; OffPeak int64 }

func TestGetRush(t *testing.T) {
	tests := map[string]struct {
		input input
		want result
	}{
		"basic offpeak": {
			input: input{ then: 10, now: 20 },
			want: result{ Rush: 0, OffPeak: 10 },
		},
		"basic rush": {
			input: input{ then: 3600*11/24, now: 3600*11/24+10 },
			want: result{ Rush: 10, OffPeak: 0 },
		},
		"long offpeak": {
			input: input{ then: 0, now: 300 },
			want: result{ Rush: 0, OffPeak: 300 },
		},
		"long rush": {
			input: input{ then: 150*11, now: 150*12 },
			want: result{ Rush: 150, OffPeak: 0 },
		},
		"over offpeak and rush": {
			input: input{ then: 0, now: 150*24 },
			want: result{ Rush: 150*5, OffPeak: 150*19 },
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rush, offPeak := GetRush(test.input.then, test.input.now)
			got := result{ Rush: rush, OffPeak: offPeak }
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetRush(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}

}

