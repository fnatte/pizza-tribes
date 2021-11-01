package internal

import (
	"time"

	. "github.com/fnatte/pizza-tribes/internal/models"
)

const ThiefSpeed = 6 * time.Minute
const ThiefCapacity = 3_000

var FullGameData = GameData{
	Buildings: map[int32]*BuildingInfo{
		int32(Building_KITCHEN): {
			Title:       "Kitchen",
			TitlePlural: "Kitchens",
			LevelInfos: []*BuildingInfo_LevelInfo{
				{
					Cost:             2_500,
					ConstructionTime: 200,
					Employer: &Employer{
						MaxWorkforce: 2,
					},
				},
				{
					Cost:             5_000,
					ConstructionTime: 450,
					Employer: &Employer{
						MaxWorkforce: 4,
					},
				},
				{
					Cost:             10_000,
					ConstructionTime: 900,
					Employer: &Employer{
						MaxWorkforce: 7,
					},
				},
				{
					Cost:             20_000,
					ConstructionTime: 3600,
					Employer: &Employer{
						MaxWorkforce: 12,
					},
				},
				{
					Cost:             40_000,
					ConstructionTime: 2 * 3600,
					Employer: &Employer{
						MaxWorkforce: 20,
					},
				},
				{
					Cost:             85_000,
					ConstructionTime: 6 * 3600,
					Employer: &Employer{
						MaxWorkforce: 35,
					},
				},
				{
					Cost:             180_000,
					ConstructionTime: 14 * 3600,
					Employer: &Employer{
						MaxWorkforce: 60,
					},
				},
				{
					Cost:             370_000,
					ConstructionTime: 24 * 3600,
					Employer: &Employer{
						MaxWorkforce: 110,
					},
				},
				{
					Cost:             750_000,
					ConstructionTime: 36 * 3600,
					Employer: &Employer{
						MaxWorkforce: 200,
					},
				},
			},
		},
		int32(Building_SHOP): {
			Title:       "Shop",
			TitlePlural: "Shops",
			LevelInfos: []*BuildingInfo_LevelInfo{
				{
					Cost:             3_500,
					ConstructionTime: 360,
					Employer: &Employer{
						MaxWorkforce: 2,
					},
				},
				{
					Cost:             10_500,
					ConstructionTime: 1 * 1200,
					Employer: &Employer{
						MaxWorkforce: 5,
					},
				},
				{
					Cost:             22_500,
					ConstructionTime: 3 * 1200,
					Employer: &Employer{
						MaxWorkforce: 9,
					},
				},
				{
					Cost:             45_000,
					ConstructionTime: 8 * 1200,
					Employer: &Employer{
						MaxWorkforce: 15,
					},
				},
				{
					Cost:             99_000,
					ConstructionTime: 14 * 3600,
					Employer: &Employer{
						MaxWorkforce: 25,
					},
				},
				{
					Cost:             200_000,
					ConstructionTime: 25 * 3600,
					Employer: &Employer{
						MaxWorkforce: 40,
					},
				},
				{
					Cost:             420_000,
					ConstructionTime: 40 * 3600,
					Employer: &Employer{
						MaxWorkforce: 70,
					},
				},
			},
		},
		int32(Building_HOUSE): {
			Title:       "House",
			TitlePlural: "Houses",
			LevelInfos: []*BuildingInfo_LevelInfo{
				{
					Cost:             3_000,
					ConstructionTime: 90,
					Residence: &Residence{
						Beds: 3,
					},
				},
				{
					Cost:             7_500,
					ConstructionTime: 200,
					Residence: &Residence{
						Beds: 6,
					},
				},
				{
					Cost:             16_000,
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
					Cost:             75_000,
					ConstructionTime: 1800,
					Residence: &Residence{
						Beds: 30,
					},
				},
				{
					Cost:             165_000,
					ConstructionTime: 4000,
					Residence: &Residence{
						Beds: 50,
					},
				},
				{
					Cost:             360_000,
					ConstructionTime: 7200,
					Residence: &Residence{
						Beds: 80,
					},
				},
				{
					Cost:             800_000,
					ConstructionTime: 2 * 7200,
					Residence: &Residence{
						Beds: 120,
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
		int32(Building_MARKETINGHQ): {
			Title:       "Marketing HQ",
			TitlePlural: "Marketing HQs",
			LevelInfos: []*BuildingInfo_LevelInfo{
				{
					Cost:             1 * 35_000,
					ConstructionTime: 1 * 3600,
					Employer: &Employer{
						MaxWorkforce: 3,
					},
				},
				{
					Cost:             2 * 35_000,
					ConstructionTime: 2 * 3600,
					Employer: &Employer{
						MaxWorkforce: 6,
					},
				},
				{
					Cost:             4 * 35_000,
					ConstructionTime: 8 * 3600,
					Employer: &Employer{
						MaxWorkforce: 12,
					},
				},
				{
					Cost:             10 * 35_000,
					ConstructionTime: 20 * 3600,
					Employer: &Employer{
						MaxWorkforce: 25,
					},
				},
				{
					Cost:             24 * 35_000,
					ConstructionTime: 48 * 3600,
					Employer: &Employer{
						MaxWorkforce: 55,
					},
				},
			},
		},
		int32(Building_RESEARCH_INSTITUTE): {
			Title:       "Research Institute",
			TitlePlural: "Research Institutes",
			LevelInfos: []*BuildingInfo_LevelInfo{
				{
					Cost:             200_000,
					ConstructionTime: 9600,
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
		int32(Education_PUBLICIST): {
			Title:       "Publicist",
			TitlePlural: "Publicists",
			Cost:        50_000,
			TrainTime:   1200,
			Employer:    Building_MARKETINGHQ.Enum(),
		},
	},
	ResearchTracks: []*ResearchTrack{
		{
			Title: "IT",
			RootNode: &ResearchNode{
				Title:        "Website",
				Discovery:    ResearchDiscovery_WEBSITE,
				Cost:         40_000,
				ResearchTime: 3600 * 2,
				Nodes: []*ResearchNode{
					{
						Title:        "Digital Ordering System",
						Discovery:    ResearchDiscovery_DIGITAL_ORDERING_SYSTEM,
						Cost:         110_000,
						ResearchTime: 3600 * 6,
						Nodes: []*ResearchNode{
							{
								Title:        "Mobile App",
								Discovery:    ResearchDiscovery_MOBILE_APP,
								Cost:         230_000,
								ResearchTime: 3600 * 24,
							},
						},
					},
				},
			},
		},
		{
			Title: "Tools",
			RootNode: &ResearchNode{
				Title:        "Masonry Oven",
				Discovery:    ResearchDiscovery_MASONRY_OVEN,
				Cost:         30_000,
				ResearchTime: 3600 * 4,
				Nodes: []*ResearchNode{
					{
						Title:        "Gas Oven",
						Discovery:    ResearchDiscovery_GAS_OVEN,
						Cost:         80_000,
						ResearchTime: 3600 * 8,
						Nodes: []*ResearchNode{
							{
								Title:        "Hybrid Oven",
								Discovery:    ResearchDiscovery_HYBRID_OVEN,
								Cost:         200_000,
								ResearchTime: 3600 * 24,
							},
						},
					},
				},
			},
		},
		{
			Title: "Pizza Craft",
			RootNode: &ResearchNode{
				Title:        "Durum Wheat",
				Discovery:    ResearchDiscovery_DURUM_WHEAT,
				Cost:         15_000,
				ResearchTime: 3600 * 2,
				Nodes: []*ResearchNode{
					{
						Title:        "Double Zero Flour",
						Discovery:    ResearchDiscovery_DOUBLE_ZERO_FLOUR,
						Cost:         150_000,
						ResearchTime: 3600 * 12,
					},
					{
						Title:        "San Marzano Tomatoes",
						Discovery:    ResearchDiscovery_SAN_MARZANO_TOMATOES,
						Cost:         120_000,
						ResearchTime: 3600 * 9,
						Nodes: []*ResearchNode{
							{
								Title:        "Ocimum Basilicum",
								Discovery:    ResearchDiscovery_OCIMUM_BASILICUM,
								Cost:         150_000,
								ResearchTime: 3600 * 10,
								Nodes: []*ResearchNode{
									{
										Title:        "Extra Virgin",
										Discovery:    ResearchDiscovery_EXTRA_VIRGIN,
										Cost:         180_000,
										ResearchTime: 3600 * 12,
									},
								},
							},
						},
					},
				},
			},
		},
	},
}
