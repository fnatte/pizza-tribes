package game

import (
	"math"
	"time"

	"github.com/fnatte/pizza-tribes/internal/game/models"
	. "github.com/fnatte/pizza-tribes/internal/game/models"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const ThiefSpeed = 6 * time.Minute
const ThiefCapacity = 3_000

var AllAppearanceParts = []*models.AppearancePart{
	{
		Id:       "redHat1",
		Category: models.AppearanceCategory_HAT,
		Free:     true,
	},
	{
		Id:       "thiefHat1",
		Category: models.AppearanceCategory_HAT,
		Free:     true,
	},
	{
		Id:       "chefHat1",
		Category: models.AppearanceCategory_HAT,
		Free:     true,
	},
	{
		Id:       "chefHat2",
		Category: models.AppearanceCategory_HAT,
	},
	{
		Id:       "guardHat1",
		Category: models.AppearanceCategory_HAT,
		Free:     true,
	},
	{
		Id:       "hat1",
		Category: models.AppearanceCategory_HAT,
	},
	{
		Id:       "hat2",
		Category: models.AppearanceCategory_HAT,
	},
	{
		Id:       "hat3",
		Category: models.AppearanceCategory_HAT,
	},
	{
		Id:       "bucketHat1",
		Category: models.AppearanceCategory_HAT,
	},
	{
		Id:       "cap1",
		Category: models.AppearanceCategory_HAT,
	},
	{
		Id:       "basicFeet1",
		Category: models.AppearanceCategory_FEET,
		Free:     true,
	},
	{
		Id:       "bigFeet1",
		Category: models.AppearanceCategory_FEET,
		Free:     true,
	},
	{
		Id:       "smallFeet1",
		Category: models.AppearanceCategory_FEET,
		Free:     true,
	},
	{
		Id:       "mixedSmile1",
		Category: models.AppearanceCategory_MOUTH,
		Free:     true,
	},
	{
		Id:       "smile1",
		Category: models.AppearanceCategory_MOUTH,
		Free:     true,
	},
	{
		Id:       "smile2",
		Category: models.AppearanceCategory_MOUTH,
		Free:     true,
	},
	{
		Id:       "smile3",
		Category: models.AppearanceCategory_MOUTH,
		Free:     true,
	},
	{
		Id:       "smile4",
		Category: models.AppearanceCategory_MOUTH,
	},
	{
		Id:       "smile5",
		Category: models.AppearanceCategory_MOUTH,
	},
	{
		Id:       "tail1",
		Category: models.AppearanceCategory_TAIL,
		Free:     true,
	},
	{
		Id:       "tail2",
		Category: models.AppearanceCategory_TAIL,
		Free:     true,
	},
	{
		Id:       "tail3",
		Category: models.AppearanceCategory_TAIL,
		Free:     true,
	},
	{
		Id:       "tail4",
		Category: models.AppearanceCategory_TAIL,
		Free:     true,
	},
	{
		Id:       "tail5",
		Category: models.AppearanceCategory_TAIL,
		Free:     true,
	},
	{
		Id:       "glasses1",
		Category: models.AppearanceCategory_EYES_EXTRA2,
		Free:     true,
	},
	{
		Id:       "glasses2",
		Category: models.AppearanceCategory_EYES_EXTRA2,
		Free:     true,
	},
	{
		Id:       "eyePatch1",
		Category: models.AppearanceCategory_EYES_EXTRA2,
		Free:     true,
	},
	{
		Id:       "monocle1",
		Category: models.AppearanceCategory_EYES_EXTRA2,
	},
	{
		Id:       "monocle2",
		Category: models.AppearanceCategory_EYES_EXTRA2,
	},
	{
		Id:       "eyeCover1",
		Category: models.AppearanceCategory_EYES_EXTRA1,
		Free:     true,
	},
	{
		Id:       "eyeStars1",
		Category: models.AppearanceCategory_EYES_EXTRA1,
		Free:     true,
	},
	{
		Id:       "eyes1",
		Category: models.AppearanceCategory_EYES,
		Free:     true,
	},
	{
		Id:       "eyes2",
		Category: models.AppearanceCategory_EYES,
		Free:     true,
	},
	{
		Id:       "eyes3",
		Category: models.AppearanceCategory_EYES,
		Free:     true,
	},
	{
		Id:       "outfit1",
		Category: models.AppearanceCategory_OUTFIT,
		Free:     true,
	},
	{
		Id:       "thiefOutfit1",
		Category: models.AppearanceCategory_OUTFIT,
		Free:     true,
	},
	{
		Id:       "guardOutfit1",
		Category: models.AppearanceCategory_OUTFIT,
		Free:     true,
	},
	{
		Id:       "mouse1",
		Category: models.AppearanceCategory_BODY,
		Free:     true,
	},
}

var AppearancePartsMap = map[string]*AppearancePart{}

var ResearchMap = map[int32]*ResearchInfo{
	// Demand
	int32(ResearchDiscovery_WEBSITE): {
		Title:       "Website",
		Description: "If only there was some kind of online medium that could increase our popularity.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Marketing",
			Value:     "+10%",
		}},
		Tree:         ResearchTree_DEMAND,
		ResearchTime: 3600 * 1,
		Requirements: []ResearchDiscovery{},
		X:            190,
		Y:            0,
	},
	int32(ResearchDiscovery_DIGITAL_ORDERING_SYSTEM): {
		Title:       "Digital Ordering System",
		Description: "With a digital ordering system the salesmice could work more effectively.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Sale speed",
			Value:     "+20%",
		}},
		Tree:         ResearchTree_DEMAND,
		ResearchTime: 3600 * 3,
		Requirements: []ResearchDiscovery{ResearchDiscovery_WEBSITE},
		X:            215,
		Y:            140,
	},
	int32(ResearchDiscovery_MOBILE_APP): {
		Title:       "Mobile App",
		Description: "A mobile app would increase our reach even further which in turn would increase demand of our fine pizzas.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Marketing",
			Value:     "+20%",
		}},
		Tree:         ResearchTree_DEMAND,
		ResearchTime: 3600 * 12,
		Requirements: []ResearchDiscovery{ResearchDiscovery_DIGITAL_ORDERING_SYSTEM},
		X:            180,
		Y:            278,
	},
	int32(ResearchDiscovery_DURUM_WHEAT): {
		Title:       "Durum Wheat",
		Description: "We should deepen our knowledge of durum wheat to improve taste of our pizzas.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Quality",
			Value:     "+5%",
		}},
		Tree:         ResearchTree_DEMAND,
		ResearchTime: 3600 * 1,
		Requirements: []ResearchDiscovery{},
		X:            50,
		Y:            20,
	},
	int32(ResearchDiscovery_DOUBLE_ZERO_FLOUR): {
		Title:       "Double Zero Flour",
		Description: "Lets continue the search for the perfect dough!",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Quality",
			Value:     "+7.5%",
		}},
		Tree:         ResearchTree_DEMAND,
		ResearchTime: 3600 * 3,
		Requirements: []ResearchDiscovery{ResearchDiscovery_DURUM_WHEAT},
		X:            24,
		Y:            153,
	},
	int32(ResearchDiscovery_SAN_MARZANO_TOMATOES): {
		Title:       "San Marzano Tomatoes",
		Description: "Our tomatoes have no taste! To improve our tomato sauce we need to find the best tomatoes.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Quality",
			Value:     "+10%",
		}},
		Tree:         ResearchTree_DEMAND,
		ResearchTime: 3600 * 4,
		Requirements: []ResearchDiscovery{ResearchDiscovery_DOUBLE_ZERO_FLOUR},
		X:            43,
		Y:            275,
	},
	int32(ResearchDiscovery_OCIMUM_BASILICUM): {
		Title:       "Ocimum Basilicum",
		Description: "A key ingredient in tomato sauce is basil. Let us learn more about the herb.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Quality",
			Value:     "+15%",
		}},
		Tree:         ResearchTree_DEMAND,
		ResearchTime: 3600 * 5,
		Requirements: []ResearchDiscovery{ResearchDiscovery_SAN_MARZANO_TOMATOES},
		X:            50,
		Y:            402,
	},
	int32(ResearchDiscovery_EXTRA_VIRGIN): {
		Title:       "Extra Virgin",
		Description: "If we could find the perfect olive oil our tomato sauce would be even tastier!",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Quality",
			Value:     "+20%",
		}},
		Tree:         ResearchTree_DEMAND,
		ResearchTime: 3600 * 6,
		Requirements: []ResearchDiscovery{ResearchDiscovery_OCIMUM_BASILICUM},
		X:            78,
		Y:            527,
	},

	// Production
	int32(ResearchDiscovery_MASONRY_OVEN): {
		Title:       "Masonry Oven",
		Description: "If we learned how to master the traditional pizza oven our pizzas would taste better &mdash; and that would lead to increased demand!",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Quality",
			Value:     "+10%",
		}},
		Tree:         ResearchTree_PRODUCTION,
		ResearchTime: 3600 * 2,
		Requirements: []ResearchDiscovery{},
		X:            50,
		Y:            20,
	},
	int32(ResearchDiscovery_GAS_OVEN): {
		Title:       "Gas Oven",
		Description: "A gas oven would heat much faster than the traditional ones. If we had gas ovens we would be able to bake pizzas faster.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Production",
			Value:     "+10%",
		}},
		Tree:         ResearchTree_PRODUCTION,
		ResearchTime: 3600 * 4,
		Requirements: []ResearchDiscovery{ResearchDiscovery_MASONRY_OVEN},
		X:            24,
		Y:            153,
	},
	int32(ResearchDiscovery_HYBRID_OVEN): {
		Title:       "Hybrid Oven",
		Description: "If we just could get the taste of traditional masonry ovens with the speed of gas ovens...",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Production",
			Value:     "+20%",
		}},
		Tree:         ResearchTree_PRODUCTION,
		ResearchTime: 3600 * 12,
		Requirements: []ResearchDiscovery{ResearchDiscovery_GAS_OVEN},
		X:            35,
		Y:            275,
	},
	int32(ResearchDiscovery_WHITEBOARD): {
		Title:       "Whiteboard",
		Description: "Wouldn't it be great with a canvas on which we could keep track of production and sales?",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Production",
			Value:     "+5%",
		}, {
			Attribute: "Sale speed",
			Value:     "+5%",
		}},
		Tree:         ResearchTree_PRODUCTION,
		ResearchTime: 3600 * 4,
		Requirements: []ResearchDiscovery{ResearchDiscovery_GAS_OVEN},
		X:            190,
		Y:            152,
	},
	int32(ResearchDiscovery_KITCHEN_STRATEGY): {
		Title:       "Kitchen Strategy",
		Description: "Let's spend some time to think about how we should organize our kitchens.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Production",
			Value:     "+15%",
		}},
		Tree:         ResearchTree_PRODUCTION,
		ResearchTime: 3600 * 7,
		Requirements: []ResearchDiscovery{ResearchDiscovery_WHITEBOARD},
		X:            200,
		Y:            289,
	},
	int32(ResearchDiscovery_STRESS_MANAGEMENT): {
		Title:       "Stress Manage­ment",
		Description: "It can get quite hectic in the pizza business. Further, it seems like mice make more mistakes when stressed -- how can we thrive in this chaotic environment?",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Production",
			Value:     "+20%",
		}, {
			Attribute: "Sale speed",
			Value:     "+20%",
		}},
		Tree:         ResearchTree_PRODUCTION,
		ResearchTime: 3600 * 12,
		Requirements: []ResearchDiscovery{ResearchDiscovery_KITCHEN_STRATEGY, ResearchDiscovery_HYBRID_OVEN},
		X:            80,
		Y:            402,
	},

	// Tapping
	int32(ResearchDiscovery_SLAM): {
		Title:       "Slam",
		Description: "Let's see if we can get more out of taps.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Tap Rewards",
			Value:     "+30%",
		}},
		Tree:         ResearchTree_TAPPING,
		ResearchTime: 1800,
		Requirements: []ResearchDiscovery{},
		X:            75,
		Y:            20,
	},
	int32(ResearchDiscovery_HIT_IT): {
		Title:       "Hit It!",
		Description: "Those instant rewards are very rewarding -- can we get even more out of it?",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Tap Rewards",
			Value:     "+50%",
		}},
		Tree:         ResearchTree_TAPPING,
		ResearchTime: 3600 * 2 + 1800,
		Requirements: []ResearchDiscovery{ResearchDiscovery_SLAM},
		X:            24,
		Y:            153,
	},
	int32(ResearchDiscovery_GRAND_SLAM): {
		Title:       "Grand Slam",
		Description: "Tap away with those juicy rewards!",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Tap Rewards",
			Value:     "+75%",
		}},
		Tree:         ResearchTree_TAPPING,
		ResearchTime: 3600 * 2,
		Requirements: []ResearchDiscovery{ResearchDiscovery_HIT_IT},
		X:            35,
		Y:            300,
	},
	int32(ResearchDiscovery_GODS_TOUCH): {
		Title:       "Gods Touch",
		Description: "Harness the greatest painting of all time -- The Creation of Jerry by Mickeylangelo.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Tap Rewards",
			Value:     "+100%",
		}},
		Tree:         ResearchTree_TAPPING,
		ResearchTime: 3600 * 5,
		Requirements: []ResearchDiscovery{ResearchDiscovery_GRAND_SLAM, ResearchDiscovery_ON_A_ROLL},
		X:            50,
		Y:            452,
	},
	int32(ResearchDiscovery_CONSECUTIVE): {
		Title:       "Consecutive",
		Description: "Those streaks are great but hard to keep alive. If only they could be more forgiving.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Prolonged streak",
			Value:     "+1h",
		}},
		Tree:         ResearchTree_TAPPING,
		ResearchTime: 3600 * 1,
		Requirements: []ResearchDiscovery{ResearchDiscovery_SLAM},
		X:            190,
		Y:            153,
	},
	int32(ResearchDiscovery_ON_A_ROLL): {
		Title:       "On a Roll",
		Description: "Unbroken series, but who keeps count?",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Prolonged streak",
			Value:     "+2h",
		}},
		Tree:         ResearchTree_TAPPING,
		ResearchTime: 3600 * 1 + 1800,
		Requirements: []ResearchDiscovery{ResearchDiscovery_CONSECUTIVE},
		X:            190,
		Y:            310,
	},

	// Thieves
	int32(ResearchDiscovery_BOOTS_OF_HASTE): {
		Title:       "Boots of Haste",
		Description: "A classic that is beloved by the veteran bandits. You wouldn't believe the struggle involved to get your hands on a pair of these back in the days. But today, it's just a research away.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Thief Speed",
			Value:     "+25%",
		}},
		Tree:         ResearchTree_THIEVES,
		ResearchTime: 3600 * 1,
		Requirements: []ResearchDiscovery{},
		X:            50,
		Y:            20,
	},
	int32(ResearchDiscovery_TIP_TOE): {
		Title:       "Tip Toe",
		Description: "Sneaky as a fox they say. But sneakier than the fox is a mouse on its toes. This reduces the risk of getting caught!",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Thief Evade",
			Value:     "+25%",
		}},
		Tree:         ResearchTree_THIEVES,
		ResearchTime: 3600 * 1 + 1800,
		Requirements: []ResearchDiscovery{ResearchDiscovery_BOOTS_OF_HASTE},
		X:            20,
		Y:            173,
	},
	int32(ResearchDiscovery_BIG_POCKETS): {
		Title:       "Big Pockets",
		Description: "It's hard to carry all those shiny coins. This allows your thieves to carry more coins.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Thief Capacity",
			Value:     "+25%",
		}},
		Tree:         ResearchTree_THIEVES,
		ResearchTime: 3600 * 1,
		Requirements: []ResearchDiscovery{},
		X:            175,
		Y:            20,
	},
	int32(ResearchDiscovery_THIEVES_FAVORITE_BAG): {
		Title:       "Thieves Faviorite Bag",
		Description: "Pockets only goes so far. Let's get real.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Thief Capacity",
			Value:     "+50%",
		}},
		Tree:         ResearchTree_THIEVES,
		ResearchTime: 3600 * 1 + 1800,
		Requirements: []ResearchDiscovery{ResearchDiscovery_BIG_POCKETS},
		X:            210,
		Y:            173,
	},
	int32(ResearchDiscovery_SHADOW_EXPERT): {
		Title:       "Shadow Expert",
		Description: "We should make darkness our ally. We weren't born in it, but we could try to eat lots of carrots.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Thief Evade",
			Value:     "+50%",
		}},
		Tree:         ResearchTree_THIEVES,
		ResearchTime: 3600 * 4,
		Requirements: []ResearchDiscovery{ResearchDiscovery_THIEVES_FAVORITE_BAG, ResearchDiscovery_TIP_TOE},
		X:            100,
		Y:            328,
	},

	// Guards
	int32(ResearchDiscovery_COFFEE): {
		Title:       "Coffee",
		Description: "A coffee bean goes a long way for a tiny mouse.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Guard Awareness",
			Value:     "+25%",
		}},
		Tree:         ResearchTree_GUARDS,
		ResearchTime: 3600 * 1,
		Requirements: []ResearchDiscovery{},
		X:            50,
		Y:            20,
	},
	int32(ResearchDiscovery_NIGHTS_WATCH): {
		Title:       "Nights Watch",
		Description: "I shall take no cheese, hold no crib, eat no seeds. I am the protector of coins and pizzas.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Guard Efficiency",
			Value:     "+15%",
		}},
		Tree:         ResearchTree_GUARDS,
		ResearchTime: 3600 * 1 + 1800,
		Requirements: []ResearchDiscovery{ResearchDiscovery_COFFEE},
		X:            20,
		Y:            173,
	},
	int32(ResearchDiscovery_TRIP_WIRE): {
		Title:       "Trip Wire",
		Description: "Simple but efficient.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Guard Efficiency",
			Value:     "+20%",
		}},
		Tree:         ResearchTree_GUARDS,
		ResearchTime: 3600 * 1,
		Requirements: []ResearchDiscovery{},
		X:            175,
		Y:            20,
	},
	int32(ResearchDiscovery_CARDIO): {
		Title:       "Cardio",
		Description: "It's not easy to pick up a chase after eating pizza. A little bit of training could do.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Guard Efficiency",
			Value:     "+25%",
		}},
		Tree:         ResearchTree_GUARDS,
		ResearchTime: 3600 * 1 + 1800,
		Requirements: []ResearchDiscovery{ResearchDiscovery_TRIP_WIRE},
		X:            210,
		Y:            173,
	},
	int32(ResearchDiscovery_LASER_ALARM): {
		Title:       "Laser Alarm",
		Description: "Let's watch some agent action movies for inspiration -- it's for research after all.",
		Rewards: []*models.ResearchInfo_Reward{{
			Attribute: "Guard Efficiency",
			Value:     "+30%",
		}},
		Tree:         ResearchTree_GUARDS,
		ResearchTime: 3600 * 3,
		Requirements: []ResearchDiscovery{ResearchDiscovery_NIGHTS_WATCH, ResearchDiscovery_CARDIO},
		X:            100,
		Y:            328,
	},
}

var FullGameData = GameData{
	AppearanceParts: AppearancePartsMap,
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
					FirstCost:             wrapperspb.Int32(0),
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
					Cost:             2_500,
					ConstructionTime: 450,
					Employer: &Employer{
						MaxWorkforce: 4,
					},
				},
				{
					Cost:             6_500,
					ConstructionTime: 900,
					Employer: &Employer{
						MaxWorkforce: 7,
					},
				},
				{
					Cost:             15_500,
					ConstructionTime: 3600,
					Employer: &Employer{
						MaxWorkforce: 12,
					},
				},
				{
					Cost:             36_500,
					ConstructionTime: 2 * 3600,
					Employer: &Employer{
						MaxWorkforce: 20,
					},
				},
				{
					Cost:             78_000,
					ConstructionTime: 3 * 3600,
					Employer: &Employer{
						MaxWorkforce: 35,
					},
				},
				{
					Cost:             160_000,
					ConstructionTime: 4 * 3600,
					Employer: &Employer{
						MaxWorkforce: 60,
					},
				},
				{
					Cost:             370_000,
					ConstructionTime: 6 * 3600,
					Employer: &Employer{
						MaxWorkforce: 110,
					},
				},
				{
					Cost:             750_000,
					ConstructionTime: 8 * 3600,
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
					FirstCost:             wrapperspb.Int32(0),
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
					Cost:             4_250,
					ConstructionTime: 1 * 1200,
					Employer: &Employer{
						MaxWorkforce: 5,
					},
				},
				{
					Cost:             16_500,
					ConstructionTime: 3 * 1200,
					Employer: &Employer{
						MaxWorkforce: 9,
					},
				},
				{
					Cost:             38_000,
					ConstructionTime: 8 * 1200,
					Employer: &Employer{
						MaxWorkforce: 15,
					},
				},
				{
					Cost:             95_000,
					ConstructionTime: 5 * 3600,
					Employer: &Employer{
						MaxWorkforce: 25,
					},
				},
				{
					Cost:             200_000,
					ConstructionTime: 8 * 3600,
					Employer: &Employer{
						MaxWorkforce: 40,
					},
				},
				{
					Cost:             420_000,
					ConstructionTime: 12 * 3600,
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
					FirstCost:             wrapperspb.Int32(100),
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
					Cost:                  30_000,
					ConstructionTime:      1500,
					FirstCost:             wrapperspb.Int32(200),
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
					Cost:             1 * 20_000,
					ConstructionTime: 1 * 3600,
					Employer: &Employer{
						MaxWorkforce: 2,
					},
				},
				{
					Cost:             2 * 25_000,
					ConstructionTime: 2 * 3600,
					Employer: &Employer{
						MaxWorkforce: 5,
					},
				},
				{
					Cost:             4 * 25_000,
					ConstructionTime: 3 * 3600,
					Employer: &Employer{
						MaxWorkforce: 10,
					},
				},
				{
					Cost:             8 * 25_000,
					ConstructionTime: 5 * 3600,
					Employer: &Employer{
						MaxWorkforce: 18,
					},
				},
				{
					Cost:             16 * 25_000,
					ConstructionTime: 7 * 3600,
					Employer: &Employer{
						MaxWorkforce: 30,
					},
				},
				{
					Cost:             32 * 25_000,
					ConstructionTime: 9 * 3600,
					Employer: &Employer{
						MaxWorkforce: 45,
					},
				},
				{
					Cost:             60 * 25_000,
					ConstructionTime: 12 * 3600,
					Employer: &Employer{
						MaxWorkforce: 65,
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
					Cost:             40_000,
					ConstructionTime: 7200,
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
				{
					Description:      "Upgrade allows sell price for pizza to changed.",
					Cost:             10_000,
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
			Cost:        8_000,
			TrainTime:   15 * 60,
			Employer:    nil,
		},
		int32(Education_THIEF): {
			Title:       "Thief",
			TitlePlural: "Thieves",
			Cost:        12_000,
			TrainTime:   30 * 60,
			Employer:    nil,
		},
		int32(Education_PUBLICIST): {
			Title:       "Publicist",
			TitlePlural: "Publicists",
			Cost:        25_000,
			TrainTime:   20 * 60,
			Employer:    Building_MARKETINGHQ.Enum(),
		},
	},
	Research: ResearchMap,
	Quests: []*models.Quest{
		{
			Id:          "1",
			Title:       "Bake and sell",
			Description: "We need to get this business going. Let's build:\n- A *Kitchen*\n- A *Shop*",
			Reward: &models.Quest_Reward{
				Coins:  550,
				Pizzas: 0,
			},
			Order: 1,
		},
		{
			Id:          "2",
			Title:       "Workforce",
			Description: "We mice to move to our town. Let's build:\n- A *House*\n- A *School*",
			Reward: &models.Quest_Reward{
				Coins:  325,
				Pizzas: 0,
			},
			Order: 2,
		},
		{
			Id:          "3",
			Title:       "Education",
			Description: "We need cooks in our kitchens and salesmice in our shops. Let's educate:\n- 1 *Chef*\n- 1 *Salesmice*",
			Reward: &models.Quest_Reward{
				Coins:  375,
				Pizzas: 50,
			},
			Order: 4,
		},
		{
			Id:          "4",
			Title:       "It takes all kinds to make a tribe",
			Description: "Your tribe is made up of individuals. Find the Town Centre, visit a mouse and change its name.",
			Reward: &models.Quest_Reward{
				OneOfItems: []string{"hat1", "cap1", "bucketHat1"},
			},
			Order: 3,
		},
		{
			Id:          "5",
			Title:       "Upgrades",
			Description: "While we could build another house, in most cases it's more efficient to upgrade your current ones.\n\nFind your house and upgrade it to level 2.",
			Reward: &models.Quest_Reward{
				Coins:  500,
				Pizzas: 0,
			},
			Order: 8,
		},
		{
			Id:          "6",
			Title:       "Knowledge",
			Description: "If you ever get stuck or need some information on the game, there's a help page. Go find it!",
			Reward: &models.Quest_Reward{
				Coins:  750,
				Pizzas: 0,
			},
			Order: 9,
		},
		{
			Id:          "7",
			Title:       "Scale up!",
			Description: "We need to ramp up our production. Let's upgrade:\n- Kitchen to level 2\n- Shop to level 2",
			Reward: &models.Quest_Reward{
				Coins:  500,
				Pizzas: 500,
			},
			Order: 10,
		},
		{
			Id:          "8",
			Title:       "Work, work",
			Description: "This tribe is getting some serious business. Let's keep going! Employ a total of 7 mice.",
			Reward: &models.Quest_Reward{
				Coins:  1_500,
				Pizzas: 300,
			},
			Order: 11,
		},
		{
			Id:          "9",
			Title:       "Statistics",
			Description: "Have you found the stats page yet? Visit the stats page and find your pizza production.\nHow many pizzas are you making per second?",
			Reward: &models.Quest_Reward{
				Coins:  1_000,
				Pizzas: 0,
			},
			Order: 12,
		},
		{
			Id:          "10",
			Title:       "Protection",
			Description: "Boss! I'm not sure we can trust our neighbors. They might be sneaking around and try to our coins! Let's educate:\n- 1 *Guard*",
			Reward: &models.Quest_Reward{
				Coins:  5_000,
				Pizzas: 1_000,
			},
			Order: 13,
		},
		{
			Id:          "11",
			Title:       "Victory is not gained by idleness",
			Description: "Chief! Looks like we have found ourselves in a competition. Go to the leaderboard and to see our progression and rank.\nWhat rank/position do we have?",
			Reward: &models.Quest_Reward{
				Coins:  1_000,
				Pizzas: 1_000,
			},
			Order: 14,
		},
		{
			Id:          "12",
			Title:       "God's touch",
			Description: "As the leader for your tribe you can tap shops and kitchens to get instant coins and pizzas. Keep your hourly streak alive to get even higher rewards!\n- Reach a tap streak of 2 (consecutive hours)",
			Reward: &models.Quest_Reward{
				Coins:  250,
				Pizzas: 250,
			},
			Order: 5,
		},
		{
			Id:          "13",
			Title:       "Supply and demand",
			Description: "By upgrading our Town Hall we can unlock the ability to change our pizza price, so that we can tune it for maximum profit!",
			Reward: &models.Quest_Reward{
				Coins:  0,
				Pizzas: 2_500,
			},
			Order: 15,
		},
		{
			Id:          "14",
			Title:       "Hats, hats, hats",
			Description: "Coins can make you win, but hats can let you do so with style.\n- Visit a mouse and change its appearance",
			Reward: &models.Quest_Reward{
				Coins:  450,
				Pizzas: 75,
			},
			Order: 6,
		},
		{
			Id:          "15",
			Title:       "Effort on Display",
			Description: "Did you visit any neighbor towns yet? If you do so you may meet their ambassador. You should select one of your own mice as ambassador to greet your friends and foes.\n- Visit a mouse and assign it as ambassador",
			Reward: &models.Quest_Reward{
				Coins:  480,
				Pizzas: 250,
			},
			Order: 7,
		},
	},
	GeniusFlashCosts: []*models.GeniusFlashCost{},
}

func polynomial1(x float64, a ...float64) float64 {
	v := 0.0
	for i := 0; i < len(a); i++ {
		v += a[i] * math.Pow(x, float64(i+1))
	}
	return v
}

func smartround(x float64) float64 {
	if x < 20_000 {
		return math.Round(x/100) * 100
	} else if x < 200_000 {
		return math.Round(x/1_000) * 1_000
	} else {
		return math.Round(x/10_000) * 10_000
	}
}

func init() {
	// Build up appearance parts map
	for _, part := range AllAppearanceParts {
		AppearancePartsMap[part.Id] = part
	}

	// Assign discovery on research
	for id, r := range ResearchMap {
		r.Discovery = models.ResearchDiscovery(id)
	}

	// Setup genius flash costs
	for n := 0; n < len(FullGameData.Research); n++ {
		cost := &GeniusFlashCost{
			Coins:  int32(smartround(polynomial1(float64(n+1), 3100, 700, -15))),
			Pizzas: int32(smartround(polynomial1(float64(n+1), 2000, 1500, -12))),
		}
		FullGameData.GeniusFlashCosts = append(FullGameData.GeniusFlashCosts, cost)
	}
}
