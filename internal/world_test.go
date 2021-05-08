package internal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type input struct{ x int; y int }
type Idx struct{ Zidx int; Eidx int }

func TestGetZoneIdx(t *testing.T) {
	tests := map[string]struct {
		input input
		want int
	}{
		"zero": {
			input: input{ x: 0, y: 0 },
			want: 0,
		},
		"x only first zone": {
			input: input{ x: 5, y: 0 },
			want: 0,
		},
		"y only first zone": {
			input: input{ x: 0, y: 5 },
			want: 0,
		},
		"second zone x": {
			input: input{ x: 15, y: 5 },
			want: 1,
		},
		"second zone y": {
			input: input{ x: 0, y: 15 },
			want: 11,
		},
		"both x and y": {
			input: input{ x: 15, y: 15 },
			want: 12,
		},
		"last zone": {
			input: input{ x: 99, y: 99 },
			want: 108,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := getZoneIdx(test.input.x, test.input.y)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetRush(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetIdx(t *testing.T) {
	tests := map[string]struct {
		input input
		want Idx
	}{
		"xy(0,0)": {
			input: input{ x: 0, y: 0 },
			want: Idx{ Zidx: 0, Eidx: 0 },
		},
		"xy(5,0)": {
			input: input{ x: 5, y: 0 },
			want: Idx{ Zidx: 0, Eidx: 5 },
		},
		"xy(0,5)": {
			input: input{ x: 0, y: 5 },
			want: Idx{ Zidx: 0, Eidx: 50 },
		},
		"xy(15,5)": {
			input: input{ x: 15, y: 5 },
			want: Idx{ Zidx: 1, Eidx: 55 },
		},
		"xy(0,15)": {
			input: input{ x: 0, y: 15 },
			want: Idx{ Zidx: 11, Eidx: 50 },
		},
		"xy(15,15)": {
			input: input{ x: 15, y: 15 },
			want: Idx{ Zidx: 12, Eidx: 55 },
		},
		"xy(99,99)": {
			input: input{ x: 99, y: 99 },
			want: Idx{ Zidx: 108, Eidx: 99 },
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			zidx, eidx := getIdx(test.input.x, test.input.y)
			if diff := cmp.Diff(test.want, Idx{ Zidx: zidx, Eidx: eidx }); diff != "" {
				t.Errorf("GetRush(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}

}
