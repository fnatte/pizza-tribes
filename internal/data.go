package internal

import (
	"time"

	. "github.com/fnatte/pizza-tribes/internal/models/gamedata"
)

func int32Ptr(i int32) *int32 {
	return &i
}

func strPtr(str string) *string {
	return &str
}

const ThiefSpeed = 6 * time.Minute
const ThiefCapacity = 3_000

var FullGameData = GameData{
	Buildings: []BuildingInfo{
		{
			ID:          "kitchen",
			Title:       "Kitchen",
			TitlePlural: "Kitchens",
			LevelInfos: []BuildingInfoLevelInfo{
				{
					Cost:             7_500,
					ConstructionTime: 200,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 2,
					},
					FirstCost:             int32Ptr(0),
					FirstConstructionTime: int32Ptr(8),
				},
				{
					Cost:             750,
					ConstructionTime: 200,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 3,
					},
				},
				{
					Cost:             3_500,
					ConstructionTime: 450,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 4,
					},
				},
				{
					Cost:             10_000,
					ConstructionTime: 900,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 7,
					},
				},
				{
					Cost:             20_000,
					ConstructionTime: 3600,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 12,
					},
				},
				{
					Cost:             40_000,
					ConstructionTime: 2 * 3600,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 20,
					},
				},
				{
					Cost:             85_000,
					ConstructionTime: 6 * 3600,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 35,
					},
				},
				{
					Cost:             180_000,
					ConstructionTime: 14 * 3600,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 60,
					},
				},
				{
					Cost:             370_000,
					ConstructionTime: 24 * 3600,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 110,
					},
				},
				{
					Cost:             750_000,
					ConstructionTime: 36 * 3600,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 200,
					},
				},
			},
		},
		{
			ID:          "shop",
			Title:       "Shop",
			TitlePlural: "Shops",
			LevelInfos: []BuildingInfoLevelInfo{
				{
					Cost:             16_500,
					ConstructionTime: 360,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 2,
					},
					FirstCost:             int32Ptr(0),
					FirstConstructionTime: int32Ptr(12),
				},
				{
					Cost:             1_000,
					ConstructionTime: 250,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 3,
					},
				},
				{
					Cost:             5_500,
					ConstructionTime: 1 * 1200,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 5,
					},
				},
				{
					Cost:             18_500,
					ConstructionTime: 3 * 1200,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 9,
					},
				},
				{
					Cost:             40_000,
					ConstructionTime: 8 * 1200,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 15,
					},
				},
				{
					Cost:             95_000,
					ConstructionTime: 14 * 3600,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 25,
					},
				},
				{
					Cost:             200_000,
					ConstructionTime: 25 * 3600,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 40,
					},
				},
				{
					Cost:             420_000,
					ConstructionTime: 40 * 3600,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 70,
					},
				},
			},
		},
		{
			ID:          "house",
			Title:       "House",
			TitlePlural: "Houses",
			LevelInfos: []BuildingInfoLevelInfo{
				{
					Cost:             8_000,
					ConstructionTime: 90,
					Residence: &BuildingInfoLevelInfoResidence{
						Beds: 3,
					},
					FirstCost:             int32Ptr(100),
					FirstConstructionTime: int32Ptr(10),
				},
				{
					Cost:             500,
					ConstructionTime: 200,
					Residence: &BuildingInfoLevelInfoResidence{
						Beds: 5,
					},
				},
				{
					Cost:             3_000,
					ConstructionTime: 200,
					Residence: &BuildingInfoLevelInfoResidence{
						Beds: 7,
					},
				},
				{
					Cost:             14_000,
					ConstructionTime: 450,
					Residence: &BuildingInfoLevelInfoResidence{
						Beds: 11,
					},
				},
				{
					Cost:             35_000,
					ConstructionTime: 900,
					Residence: &BuildingInfoLevelInfoResidence{
						Beds: 18,
					},
				},
				{
					Cost:             75_000,
					ConstructionTime: 1800,
					Residence: &BuildingInfoLevelInfoResidence{
						Beds: 30,
					},
				},
				{
					Cost:             165_000,
					ConstructionTime: 4000,
					Residence: &BuildingInfoLevelInfoResidence{
						Beds: 50,
					},
				},
				{
					Cost:             360_000,
					ConstructionTime: 7200,
					Residence: &BuildingInfoLevelInfoResidence{
						Beds: 80,
					},
				},
				{
					Cost:             800_000,
					ConstructionTime: 2 * 7200,
					Residence: &BuildingInfoLevelInfoResidence{
						Beds: 120,
					},
				},
			},
		},
		{
			ID:          "school",
			Title:       "School",
			TitlePlural: "Schools",
			MaxCount:    int32Ptr(1),
			LevelInfos: []BuildingInfoLevelInfo{
				{
					Cost:                  30_000,
					ConstructionTime:      1500,
					FirstCost:             int32Ptr(200),
					FirstConstructionTime: int32Ptr(15),
				},
			},
		},
		{
			ID:          "marketinghq",
			Title:       "Marketing HQ",
			TitlePlural: "Marketing HQs",
			MaxCount:    int32Ptr(1),
			LevelInfos: []BuildingInfoLevelInfo{
				{
					Cost:             1 * 30_000,
					ConstructionTime: 1 * 3600,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 3,
					},
				},
				{
					Cost:             2 * 30_000,
					ConstructionTime: 2 * 3600,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 6,
					},
				},
				{
					Cost:             5 * 30_000,
					ConstructionTime: 8 * 3600,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 12,
					},
				},
				{
					Cost:             12 * 30_000,
					ConstructionTime: 20 * 3600,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 25,
					},
				},
				{
					Cost:             30 * 30_000,
					ConstructionTime: 48 * 3600,
					Employer: &BuildingInfoLevelInfoEmployer{
						MaxWorkforce: 55,
					},
				},
			},
		},
		{
			ID:          "research_institute",
			Title:       "Research Institute",
			TitlePlural: "Research Institutes",
			MaxCount:    int32Ptr(1),
			LevelInfos: []BuildingInfoLevelInfo{
				{
					Cost:             200_000,
					ConstructionTime: 9600,
				},
			},
		},
		{
			ID:          "town_centre",
			Title:       "Town Centre",
			TitlePlural: "Town Centres",
			MaxCount:    int32Ptr(1),
			LevelInfos: []BuildingInfoLevelInfo{
				{
					Cost:             0,
					ConstructionTime: 0,
				},
			},
		},
	},
	Educations: []EducationInfo{
		{
			ID:          "chef",
			Title:       "Chef",
			TitlePlural: "Chefs",
			Cost:        0,
			TrainTime:   200,
			Employer:    strPtr("kitchen"),
		},
		{
			ID:          "salesmouse",
			Title:       "Salesmouse",
			TitlePlural: "Salesmice",
			Cost:        0,
			TrainTime:   100,
			Employer:    strPtr("shop"),
		},
		{
			ID:          "guard",
			Title:       "Security Guard",
			TitlePlural: "Security Guards",
			Cost:        10_000,
			TrainTime:   1000,
			Employer:    nil,
		},
		{
			ID:          "thief",
			Title:       "Thief",
			TitlePlural: "Thieves",
			Cost:        20_000,
			TrainTime:   1800,
			Employer:    nil,
		},
		{
			ID:          "publicist",
			Title:       "Publicist",
			TitlePlural: "Publicists",
			Cost:        50_000,
			TrainTime:   1200,
			Employer:    strPtr("marketinghq"),
		},
	},
	ResearchTracks: []ResearchTrack{
		{
			Title: "IT",
			RootNode: ResearchNode{
				Title:        "Website",
				Discovery:    "website",
				Cost:         40_000,
				ResearchTime: 3600 * 2,
				Nodes: []ResearchNode{
					{
						Title:        "Digital Ordering System",
						Discovery:    "digital_ordering_system",
						Cost:         110_000,
						ResearchTime: 3600 * 6,
						Nodes: []ResearchNode{
							{
								Title:        "Mobile App",
								Discovery:    "mobile_app",
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
			RootNode: ResearchNode{
				Title:        "Masonry Oven",
				Discovery:    "masonry_oven",
				Cost:         30_000,
				ResearchTime: 3600 * 4,
				Nodes: []ResearchNode{
					{
						Title:        "Gas Oven",
						Discovery:    "gas_oven",
						Cost:         80_000,
						ResearchTime: 3600 * 8,
						Nodes: []ResearchNode{
							{
								Title:        "Hybrid Oven",
								Discovery:    "hybrid_oven",
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
			RootNode: ResearchNode{
				Title:        "Durum Wheat",
				Discovery:    "durum_wheat",
				Cost:         15_000,
				ResearchTime: 3600 * 2,
				Nodes: []ResearchNode{
					{
						Title:        "Double Zero Flour",
						Discovery:    "double_zero_flour",
						Cost:         150_000,
						ResearchTime: 3600 * 12,
					},
					{
						Title:        "San Marzano Tomatoes",
						Discovery:    "san_marzano_tomatoes",
						Cost:         120_000,
						ResearchTime: 3600 * 9,
						Nodes: []ResearchNode{
							{
								Title:        "Ocimum Basilicum",
								Discovery:    "ocimum_basilicum",
								Cost:         150_000,
								ResearchTime: 3600 * 10,
								Nodes: []ResearchNode{
									{
										Title:        "Extra Virgin",
										Discovery:    "extra_virgin",
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
	Quests: []Quest{
		{
			ID:          "1",
			Title:       "Bake and sell",
			Description: "We need to get this business going. Let's build:\n- A *Kitchen*\n- A *Shop*",
			Reward: QuestReward{
				Coins:  550,
				Pizzas: 0,
			},
		},
		{
			ID:          "2",
			Title:       "Workforce",
			Description: "We mice to move to our town. Let's build:\n- A *House*\n- A *School*",
			Reward: QuestReward{
				Coins:  325,
				Pizzas: 0,
			},
		},
		{
			ID:          "3",
			Title:       "Education",
			Description: "We need cooks in our kitchens and salesmice in our shops. Let's educate:\n- 1 *Chef*\n- 1 *Salesmice*",
			Reward: QuestReward{
				Coins:  375,
				Pizzas: 50,
			},
		},
		{
			ID:          "4",
			Title:       "It takes all kinds to make a tribe",
			Description: "Your tribe is made up of individuals. Find the Town Centre, visit a mouse and change its name.",
			Reward: QuestReward{
				Coins:  450,
				Pizzas: 75,
			},
		},
		{
			ID:          "5",
			Title:       "Upgrades",
			Description: "While we could build another house, in most cases it's more efficient to upgrade your current ones.\n\nFind your house and upgrade it to level 2.",
			Reward: QuestReward{
				Coins:  500,
				Pizzas: 0,
			},
		},
		{
			ID:          "6",
			Title:       "Knowledge",
			Description: "If you ever get stuck or need some information on the game, there's a help page. Go find it!",
			Reward: QuestReward{
				Coins:  750,
				Pizzas: 0,
			},
		},
		{
			ID:          "7",
			Title:       "Scale up!",
			Description: "We need to ramp up our production. Let's upgrade:\n- Kitchen to level 2\n- Shop to level 2",
			Reward: QuestReward{
				Coins:  500,
				Pizzas: 500,
			},
		},
	},
}
