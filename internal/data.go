package internal

var FullGameData = GameData{
	Buildings: map[int32]*BuildingInfo{
		int32(Building_KITCHEN): {
			Title: "Kitchen",
			TitlePlural: "Kitchens",
			Cost: 250,
			ConstructionTime: 900,
			Employer: &Employer{
				MaxWorkforce: 7,
			},
		},
		int32(Building_SHOP): {
			Title: "Shop",
			TitlePlural: "Shops",
			Cost: 750,
			ConstructionTime: 1200,
			Employer: &Employer{
				MaxWorkforce: 5,
			},
		},
		int32(Building_HOUSE): {
			Title: "House",
			TitlePlural: "Houses",
			Cost: 1_000,
			ConstructionTime: 450,
			Employer: nil,
		},
		int32(Building_SCHOOL): {
			Title: "School",
			TitlePlural: "Schools",
			Cost: 10_000,
			ConstructionTime: 3600,
			Employer: nil,
		},
	},
	Educations: map[int32]*EducationInfo{
		int32(Education_CHEF): {
			Title: "Chef",
			TitlePlural: "Chefs",
			Cost: 100,
			TrainTime: 200,
			Employer: Building_KITCHEN.Enum(),
		},
		int32(Education_SALESMOUSE): {
			Title: "Salesmouse",
			TitlePlural: "Salesmice",
			Cost: 150,
			TrainTime: 100,
			Employer: Building_SHOP.Enum(),
		},
		int32(Education_GUARD): {
			Title: "Guard",
			TitlePlural: "Guards",
			Cost: 500,
			TrainTime: 1000,
			Employer: nil,
		},
		int32(Education_THIEF): {
			Title: "Thief",
			TitlePlural: "Thieves",
			Cost: 1500,
			TrainTime: 1800,
			Employer: nil,
		},
	},
}
