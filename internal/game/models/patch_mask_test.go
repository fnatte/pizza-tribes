package models

import (
	"testing"
	"unicode"

	"github.com/google/go-cmp/cmp"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

func TestGetValueByPath(t *testing.T) {
	gs := &GameState{}
	fd := gs.ProtoReflect().Descriptor().Fields().
		ByName(protoreflect.Name("constructionQueue"))
	gs.ProtoReflect().Get(fd)

	tests := map[string]struct {
		path string
		want interface{}
	}{
		"simple": {
			path: "townX",
			want: int32(10),
		},
		"nested inside map": {
			path: "lots.2.level",
			want: int32(2),
		},
		"map value": {
			path: "lots.ab",
			want: &GameState_Lot{
				Level: 4,
				Taps:  5,
			},
		},
		"map": {
			path: "lots",
			want: map[string]*GameState_Lot{
				"2": {
					Level: 2,
					Taps:  3,
				},
				"ab": {
					Level: 4,
					Taps:  5,
				},
			},
		},
		"map nil value": {
			path: "lots.abc",
			want: nil,
		},
		"nested in message": {
			path: "resources.coins",
			want: int32(10),
		},
		"nested value at first position in list": {
			path: "travelQueue.0.destinationX",
			want: int32(2),
		},
		"nested value at second position in list": {
			path: "travelQueue.1.destinationX",
			want: int32(4),
		},
		"slice": {
			path: "constructionQueue",
			want: []*Construction{
				{
					LotId: "2",
					Level: 3,
				},
			},
		},
		"empty slice": {
			path: "researchQueue",
			want: []*OngoingResearch{},
		},
		"enum slice": {
			path: "discoveries",
			want: []ResearchDiscovery{ResearchDiscovery_DURUM_WHEAT, ResearchDiscovery_MOBILE_APP},
		},
		"string slice": {
			path: "appearanceParts",
			want: []string{"test"},
		},
		// TODO: add case for invalid field
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gs := &GameState{
				TownX: 10,
				TownY: 12,
				Resources: &GameState_Resources{
					Coins: 10,
				},
				Lots: map[string]*GameState_Lot{
					"2": {
						Level: 2,
						Taps:  3,
					},
					"ab": {
						Level: 4,
						Taps:  5,
					},
				},
				TravelQueue: []*Travel{
					{
						DestinationX: 2,
						DestinationY: 3,
					},
					{
						DestinationX: 4,
						DestinationY: 5,
					},
				},
				ConstructionQueue: []*Construction{
					{
						LotId: "2",
						Level: 3,
					},
				},
				ResearchQueue: []*OngoingResearch{},
				AppearanceParts: []string{"test"},
				Discoveries: []ResearchDiscovery{ResearchDiscovery_DURUM_WHEAT, ResearchDiscovery_MOBILE_APP},
			}

			got, err := GetValueByPath(gs, test.path)
			if err != nil {
				t.Errorf("GetValueByPath returned an error: %v", err)
			}

			opt := cmp.FilterPath(func(p cmp.Path) bool {
				ps := p.String()
				return len(ps) > 0 && unicode.IsLower([]rune(ps)[0])
			}, cmp.Ignore())

			if diff := cmp.Diff(test.want, got, opt); diff != "" {
				t.Errorf("GetValueByPath(...) mismatch (-want +got):\n%s", diff)
			}
		})
	}

}
