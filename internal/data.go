package internal

import (
	"time"
	. "github.com/fnatte/pizza-tribes/internal/models"
)

const ThiefSpeed = 5 * time.Minute
const ThiefCapacity = 4_000

var FullGameData = GameData{
	Buildings: map[int32]*BuildingInfo{
		int32(Building_KITCHEN): {
			Title:       "Kitchen",
			TitlePlural: "Kitchens",
			LevelInfos: []*BuildingInfo_LevelInfo{
				{
					Cost:             10_000,
					ConstructionTime: 900,
					Employer: &Employer{
						MaxWorkforce: 7,
					},
				},
				{
					Cost:             20_000,
					ConstructionTime: 1800,
					Employer: &Employer{
						MaxWorkforce: 12,
					},
				},
				{
					Cost:             40_000,
					ConstructionTime: 2700,
					Employer: &Employer{
						MaxWorkforce: 20,
					},
				},
				{
					Cost:             85_000,
					ConstructionTime: 4800,
					Employer: &Employer{
						MaxWorkforce: 35,
					},
				},
			},
		},
		int32(Building_SHOP): {
			Title:       "Shop",
			TitlePlural: "Shops",
			LevelInfos: []*BuildingInfo_LevelInfo{
				{
					Cost:             12_500,
					ConstructionTime: 1200,
					Employer: &Employer{
						MaxWorkforce: 5,
					},
				},
				{
					Cost:             25_000,
					ConstructionTime: 2400,
					Employer: &Employer{
						MaxWorkforce: 9,
					},
				},
				{
					Cost:             48_000,
					ConstructionTime: 3600,
					Employer: &Employer{
						MaxWorkforce: 15,
					},
				},
				{
					Cost:             99_000,
					ConstructionTime: 8000,
					Employer: &Employer{
						MaxWorkforce: 25,
					},
				},
			},
		},
		int32(Building_HOUSE): {
			Title:       "House",
			TitlePlural: "Houses",
			LevelInfos: []*BuildingInfo_LevelInfo{
				{
					Cost:             17_000,
					ConstructionTime: 450,
					Residence: &Residence{
						Beds: 10,
					},
				},
				{
					Cost:             35_000,
					ConstructionTime: 900,
					Residence: &Residence{
						Beds: 18,
					},
				},
				{
					Cost:             60_000,
					ConstructionTime: 1800,
					Residence: &Residence{
						Beds: 30,
					},
				},
				{
					Cost:             135_000,
					ConstructionTime: 4000,
					Residence: &Residence{
						Beds: 50,
					},
				},
				{
					Cost:             250_000,
					ConstructionTime: 7200,
					Residence: &Residence{
						Beds: 90,
					},
				},
			},
		},
		int32(Building_SCHOOL): {
			Title:       "School",
			TitlePlural: "Schools",
			LevelInfos: []*BuildingInfo_LevelInfo{
				{
					Cost:             30_000,
					ConstructionTime: 3600,
				},
			},
		},
	},
	Educations: map[int32]*EducationInfo{
		int32(Education_CHEF): {
			Title:       "Chef",
			TitlePlural: "Chefs",
			Cost:        0,
			TrainTime:   200,
			Employer:    Building_KITCHEN.Enum(),
		},
		int32(Education_SALESMOUSE): {
			Title:       "Salesmouse",
			TitlePlural: "Salesmice",
			Cost:        0,
			TrainTime:   100,
			Employer:    Building_SHOP.Enum(),
		},
		int32(Education_GUARD): {
			Title:       "Security Guard",
			TitlePlural: "Security Guards",
			Cost:        10_000,
			TrainTime:   1000,
			Employer:    nil,
		},
		int32(Education_THIEF): {
			Title:       "Thief",
			TitlePlural: "Thieves",
			Cost:        20_000,
			TrainTime:   1800,
			Employer:    nil,
		},
	},
}
