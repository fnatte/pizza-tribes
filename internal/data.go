package internal

import (
	"time"

	"github.com/fnatte/pizza-tribes/internal/models"
	. "github.com/fnatte/pizza-tribes/internal/models"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
					Cost:             7_500,
					ConstructionTime: 200,
					Employer: &Employer{
						MaxWorkforce: 2,
					},
					FirstCost: wrapperspb.Int32(0),
					FirstConstructionTime: wrapperspb.Int32(8),
				},
				{
					Cost:             750,
					ConstructionTime: 200,
					Employer: &Employer{
						MaxWorkforce: 3,
					},
				},
				{
					Cost:             3_500,
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
					Cost:             16_500,
					ConstructionTime: 360,
					Employer: &Employer{
						MaxWorkforce: 2,
					},
					FirstCost: wrapperspb.Int32(0),
					FirstConstructionTime: wrapperspb.Int32(12),
				},
				{
					Cost:             1_000,
					ConstructionTime: 250,
					Employer: &Employer{
						MaxWorkforce: 3,
					},
				},
				{
					Cost:             5_500,
					ConstructionTime: 1 * 1200,
					Employer: &Employer{
						MaxWorkforce: 5,
					},
				},
				{
					Cost:             18_500,
					ConstructionTime: 3 * 1200,
					Employer: &Employer{
						MaxWorkforce: 9,
					},
				},
				{
					Cost:             40_000,
					ConstructionTime: 8 * 1200,
					Employer: &Employer{
						MaxWorkforce: 15,
					},
				},
				{
					Cost:             95_000,
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
					Cost:             8_000,
					ConstructionTime: 30,
					Residence: &Residence{
						Beds: 3,
					},
					FirstCost: wrapperspb.Int32(100),
					FirstConstructionTime: wrapperspb.Int32(10),
				},
				{
					Cost:             500,
					ConstructionTime: 60,
					Residence: &Residence{
						Beds: 5,
					},
				},
				{
					Cost:             3_000,
					ConstructionTime: 200,
					Residence: &Residence{
						Beds: 7,
					},
				},
				{
					Cost:             14_000,
					ConstructionTime: 450,
					Residence: &Residence{
						Beds: 11,
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
			MaxCount:    wrapperspb.Int32(1),
			LevelInfos: []*BuildingInfo_LevelInfo{
				{
					Cost:             30_000,
					ConstructionTime: 1500,
					FirstCost: wrapperspb.Int32(200),
					FirstConstructionTime: wrapperspb.Int32(15),
				},
			},
		},
		int32(Building_MARKETINGHQ): {
			Title:       "Marketing HQ",
			TitlePlural: "Marketing HQs",
			MaxCount:    wrapperspb.Int32(1),
			LevelInfos: []*BuildingInfo_LevelInfo{
				{
					Cost:             1 * 30_000,
					ConstructionTime: 1 * 3600,
					Employer: &Employer{
						MaxWorkforce: 3,
					},
				},
				{
					Cost:             2 * 30_000,
					ConstructionTime: 2 * 3600,
					Employer: &Employer{
						MaxWorkforce: 6,
					},
				},
				{
					Cost:             5 * 30_000,
					ConstructionTime: 8 * 3600,
					Employer: &Employer{
						MaxWorkforce: 12,
					},
				},
				{
					Cost:             12 * 30_000,
					ConstructionTime: 20 * 3600,
					Employer: &Employer{
						MaxWorkforce: 25,
					},
				},
				{
					Cost:             30 * 30_000,
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
			MaxCount:    wrapperspb.Int32(1),
			LevelInfos: []*BuildingInfo_LevelInfo{
				{
					Cost:             200_000,
					ConstructionTime: 9600,
				},
			},
		},
		int32(Building_TOWN_CENTRE): {
			Title:       "Town Centre",
			TitlePlural: "Town Centres",
			MaxCount:    wrapperspb.Int32(1),
			LevelInfos: []*BuildingInfo_LevelInfo{
				{
					Cost:             0,
					ConstructionTime: 0,
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
	Quests: []*models.Quest{
		{
			Id:          "1",
			Title:       "Bake and sell",
			Description: "We need to get this business going. Let's build:\n- A *Kitchen*\n- A *Shop*",
			Reward: &models.Quest_Reward{
				Coins:  550,
				Pizzas: 0,
			},
		},
		{
			Id:          "2",
			Title:       "Workforce",
			Description: "We mice to move to our town. Let's build:\n- A *House*\n- A *School*",
			Reward: &models.Quest_Reward{
				Coins:  325,
				Pizzas: 0,
			},
		},
		{
			Id:          "3",
			Title:       "Education",
			Description: "We need cooks in our kitchens and salesmice in our shops. Let's educate:\n- 1 *Chef*\n- 1 *Salesmice*",
			Reward: &models.Quest_Reward{
				Coins:  375,
				Pizzas: 50,
			},
		},
		{
			Id:          "4",
			Title:       "It takes all kinds to make a tribe",
			Description: "Your tribe is made up of individuals. Find the Town Centre, visit a mouse and change its name.",
			Reward: &models.Quest_Reward{
				Coins:  450,
				Pizzas: 75,
			},
		},
		{
			Id:          "5",
			Title:       "Upgrades",
			Description: "While we could build another house, in most cases it's more efficient to upgrade your current ones.\n\nFind your house and upgrade it to level 2.",
			Reward: &models.Quest_Reward{
				Coins:  500,
				Pizzas: 0,
			},
		},
		{
			Id:          "6",
			Title:       "Knowledge",
			Description: "If you ever get stuck or need some information on the game, there's a help page. Go find it!",
			Reward: &models.Quest_Reward{
				Coins:  750,
				Pizzas: 0,
			},
		},
		{
			Id:          "7",
			Title:       "Scale up!",
			Description: "We need to ramp up our production. Let's upgrade:\n- Kitchen to level 2\n- Shop to level 2",
			Reward: &models.Quest_Reward{
				Coins:  500,
				Pizzas: 500,
			},
		},
	},
}
