package internal

import "time"

const MicePerHouse = 10

const ThiefSpeed = 2 * time.Minute

var FullGameData = GameData{
	Buildings: map[int32]*BuildingInfo{
		int32(Building_KITCHEN): {
			Title: "Kitchen",
			TitlePlural: "Kitchens",
			Cost: 10_000,
			ConstructionTime: 900,
			Employer: &Employer{
				MaxWorkforce: 7,
			},
		},
		int32(Building_SHOP): {
			Title: "Shop",
			TitlePlural: "Shops",
			Cost: 12_500,
			ConstructionTime: 1200,
			Employer: &Employer{
				MaxWorkforce: 5,
			},
		},
		int32(Building_HOUSE): {
			Title: "House",
			TitlePlural: "Houses",
			Cost: 17_000,
			ConstructionTime: 450,
			Employer: nil,
		},
		int32(Building_SCHOOL): {
			Title: "School",
			TitlePlural: "Schools",
			Cost: 30_000,
			ConstructionTime: 3600,
			Employer: nil,
		},
	},
	Educations: map[int32]*EducationInfo{
		int32(Education_CHEF): {
			Title: "Chef",
			TitlePlural: "Chefs",
			Cost: 0,
			TrainTime: 200,
			Employer: Building_KITCHEN.Enum(),
		},
		int32(Education_SALESMOUSE): {
			Title: "Salesmouse",
			TitlePlural: "Salesmice",
			Cost: 0,
			TrainTime: 100,
			Employer: Building_SHOP.Enum(),
		},
		int32(Education_GUARD): {
			Title: "Guard",
			TitlePlural: "Guards",
			Cost: 10_000,
			TrainTime: 1000,
			Employer: nil,
		},
		int32(Education_THIEF): {
			Title: "Thief",
			TitlePlural: "Thieves",
			Cost: 20_000,
			TrainTime: 1800,
			Employer: nil,
		},
	},
}
