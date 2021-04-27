package internal

var FullGameData = GameData{
	BuildingInfos: map[int32]*BuildingInfo{
		int32(Building_KITCHEN): {
			Employer: &Employer{
				MaxWorkforce: 7,
			},
		},
		int32(Building_SHOP): {
			Employer: &Employer{
				MaxWorkforce: 5,
			},
		},
		int32(Building_HOUSE): {
			Employer: nil,
		},
		int32(Building_SCHOOL): {
			Employer: nil,
		},
	},
}
