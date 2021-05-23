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
		int32(Building_MARKETINGHQ): {
			Title:       "Marketing HQ",
			TitlePlural: "Marketing HQs",
			LevelInfos: []*BuildingInfo_LevelInfo{
				{
					Cost:             100_000,
					ConstructionTime: 7200,
					Employer: &Employer{
						MaxWorkforce: 5,
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
			Cost:        80_000,
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
				Cost:         50_000,
				ResearchTime: 3600 * 2,
				Nodes: []*ResearchNode{
					{
						Title:        "Digital Ordering System",
						Discovery:    ResearchDiscovery_DIGITAL_ORDERING_SYSTEM,
						Cost:         125_000,
						ResearchTime: 3600 * 6,
						Nodes: []*ResearchNode{
							{
								Title:        "Mobile App",
								Discovery:    ResearchDiscovery_MOBILE_APP,
								Cost:         250_000,
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
				Cost:         40_000,
				ResearchTime: 3600 * 4,
				Nodes: []*ResearchNode{
					{
						Title:        "Gas Oven",
						Discovery:    ResearchDiscovery_GAS_OVEN,
						Cost:         100_000,
						ResearchTime: 3600 * 8,
						Nodes: []*ResearchNode{
							{
								Title:        "Hybrid Oven",
								Discovery:    ResearchDiscovery_HYBRID_OVEN,
								Cost:         250_000,
								ResearchTime: 3600 * 24,
							},
						},
					},
				},
			},
		},
		{
			Title: "Pizza",
			RootNode: &ResearchNode{
				Title:        "Durum Wheat",
				Discovery:    ResearchDiscovery_DURUM_WHEAT,
				Cost:         20_000,
				ResearchTime: 3600 * 2,
				Nodes: []*ResearchNode{
					{
						Title:        "Double Zero Flour",
						Discovery:    ResearchDiscovery_DOUBLE_ZERO_FLOUR,
						Cost:         180_000,
						ResearchTime: 3600 * 12,
					},
					{
						Title:        "San Marzano Tomatoes",
						Discovery:    ResearchDiscovery_SAN_MARZANO_TOMATOES,
						Cost:         150_000,
						ResearchTime: 3600 * 9,
						Nodes: []*ResearchNode{
							{
								Title:        "Ocimum Basilicum",
								Discovery:    ResearchDiscovery_OCIMUM_BASILICUM,
								Cost:         180_000,
								ResearchTime: 3600 * 10,
								Nodes: []*ResearchNode{
									{
										Title:        "Extra Virgin",
										Discovery:    ResearchDiscovery_EXTRA_VIRGIN,
										Cost:         200_000,
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
