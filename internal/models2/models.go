// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    building, err := UnmarshalBuilding(bytes)
//    bytes, err = building.Marshal()
//
//    education, err := UnmarshalEducation(bytes)
//    bytes, err = education.Marshal()
//
//    gameData, err := UnmarshalGameData(bytes)
//    bytes, err = gameData.Marshal()
//
//    gameState, err := UnmarshalGameState(bytes)
//    bytes, err = gameState.Marshal()
//
//    researchDiscovery, err := UnmarshalResearchDiscovery(bytes)
//    bytes, err = researchDiscovery.Marshal()

package models2

import "encoding/json"

func UnmarshalBuilding(data []byte) (Building, error) {
	var r Building
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Building) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalEducation(data []byte) (Education, error) {
	var r Education
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Education) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalGameData(data []byte) (GameData, error) {
	var r GameData
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GameData) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalGameState(data []byte) (GameState, error) {
	var r GameState
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GameState) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalResearchDiscovery(data []byte) (ResearchDiscovery, error) {
	var r ResearchDiscovery
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ResearchDiscovery) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GameData struct {
	Buildings      []BuildingInfo  `json:"buildings"`               
	Educations     []EducationInfo `json:"educations"`              
	Quests         []QuestInfo     `json:"quests"`                  
	ResearchTracks []ResearchTrack `json:"researchTracks,omitempty"`
}

type BuildingInfo struct {
	ID          string      `json:"id"`         
	LevelInfos  []LevelInfo `json:"levelInfos"` 
	MaxCount    *int64      `json:"maxCount"`   
	Title       string      `json:"title"`      
	TitlePlural string      `json:"titlePlural"`
}

type LevelInfo struct {
	ConstructionTime      int64      `json:"constructionTime"`     
	Cost                  int64      `json:"cost"`                 
	Employer              *Employer  `json:"employer,omitempty"`   
	FirstConstructionTime *int64     `json:"firstConstructionTime"`
	FirstCost             *int64     `json:"firstCost"`            
	Residence             *Residence `json:"residence,omitempty"`  
}

type Employer struct {
	MaxWorkforce int64 `json:"maxWorkforce"`
}

type Residence struct {
	Beds int64 `json:"beds"`
}

type EducationInfo struct {
	Cost        int64   `json:"cost"`       
	Employer    *string `json:"employer"`   
	ID          string  `json:"id"`         
	Title       string  `json:"title"`      
	TitlePlural string  `json:"titlePlural"`
	TrainTime   int64   `json:"trainTime"`  
}

type QuestInfo struct {
	Description string      `json:"description"`
	ID          string      `json:"id"`         
	Reward      QuestReward `json:"reward"`     
	Title       string      `json:"title"`      
}

type QuestReward struct {
	Coins  int64 `json:"coins"` 
	Pizzas int64 `json:"pizzas"`
}

type ResearchTrack struct {
	RootNode ResearchNode `json:"rootNode"`
	Title    string       `json:"title"`   
}

type ResearchNode struct {
	Cost         int64             `json:"cost"`           
	Discovery    ResearchDiscovery `json:"discovery"`      
	Nodes        []ResearchNode    `json:"nodes,omitempty"`
	ResearchTime int64             `json:"research_time"`  
	Title        string            `json:"title"`          
}

type GameState struct {
	ConstructionQueue []ConstructionQueue   `json:"constructionQueue"`
	Discoveries       []ResearchDiscovery   `json:"discoveries"`      
	Lots              map[string]Lot        `json:"lots"`             
	Mice              map[string]Mouse      `json:"mice"`             
	Population        Population            `json:"population"`       
	Quests            map[string]QuestState `json:"quests"`           
	ResearchQueue     []ResearchQueue       `json:"researchQueue"`    
	Resources         Resources             `json:"resources"`        
	Stats             Stats                 `json:"stats"`            
	Timestamp         string                `json:"timestamp"`        
	TownX             int64                 `json:"townX"`            
	TownY             int64                 `json:"townY"`            
	TrainingQueue     []Training            `json:"trainingQueue"`    
	TravelQueue       []Travel              `json:"travelQueue"`      
}

type ConstructionQueue struct {
	Building   Building `json:"building"`  
	CompleteAt string   `json:"completeAt"`
	Level      int64    `json:"level"`     
	LotID      string   `json:"lotId"`     
	Razing     bool     `json:"razing"`    
}

type Lot struct {
	Building Building `json:"building"`
	Level    int64    `json:"level"`   
	Streak   int64    `json:"streak"`  
	TappedAt int64    `json:"tappedAt"`
	Taps     int64    `json:"taps"`    
}

type Mouse struct {
	Education       *Education `json:"education,omitempty"`
	IsBeingEducated bool       `json:"isBeingEducated"`    
	IsEducated      bool       `json:"isEducated"`         
	Name            string     `json:"name"`               
}

type Population struct {
	Chefs      int64 `json:"chefs"`     
	Guards     int64 `json:"guards"`    
	Publicists int64 `json:"publicists"`
	Salesmice  int64 `json:"salesmice"` 
	Thieves    int64 `json:"thieves"`   
	Uneducated int64 `json:"uneducated"`
}

type QuestState struct {
	ClaimedReward bool `json:"claimedReward"`
	Completed     bool `json:"completed"`    
	Opened        bool `json:"opened"`       
}

type ResearchQueue struct {
	CompleteAt string             `json:"completeAt"`        
	Dicovery   *ResearchDiscovery `json:"dicovery,omitempty"`
}

type Resources struct {
	Coins  int64 `json:"coins"` 
	Pizzas int64 `json:"pizzas"`
}

type Stats struct {
	DemandOffpeak           float64 `json:"demandOffpeak"`          
	DemandRushHour          float64 `json:"demandRushHour"`         
	EmployedChefs           int64   `json:"employedChefs"`          
	EmployedSalesmice       int64   `json:"employedSalesmice"`      
	MaxSellsByMicePerSecond float64 `json:"maxSellsByMicePerSecond"`
	PizzasProducedPerSecond float64 `json:"pizzasProducedPerSecond"`
}

type Training struct {
	Amount     int64      `json:"amount"`             
	CompleteAt string     `json:"completeAt"`         
	Education  *Education `json:"education,omitempty"`
}

type Travel struct {
	ArrivalAt    string `json:"arrivalAt"`   
	Coins        string `json:"coins"`       
	DestinationX int64  `json:"destinationX"`
	DestinationY int64  `json:"destinationY"`
	Returning    bool   `json:"returning"`   
	Thieves      int64  `json:"thieves"`     
}

type ResearchDiscovery string
const (
	DigitalOrderingSystem ResearchDiscovery = "digital_ordering_system"
	DoubleZeroFlour ResearchDiscovery = "double_zero_flour"
	DurumWheat ResearchDiscovery = "durum_wheat"
	ExtraVirgin ResearchDiscovery = "extra_virgin"
	GasOven ResearchDiscovery = "gas_oven"
	HybridOven ResearchDiscovery = "hybrid_oven"
	MasonryOven ResearchDiscovery = "masonry_oven"
	MobileApp ResearchDiscovery = "mobile_app"
	OcimumBasilicum ResearchDiscovery = "ocimum_basilicum"
	SANMarzanoTomatoes ResearchDiscovery = "san_marzano_tomatoes"
	Website ResearchDiscovery = "website"
)

type Building string
const (
	House Building = "house"
	Kitchen Building = "kitchen"
	Marketinghq Building = "marketinghq"
	ResearchInstitute Building = "research_institute"
	School Building = "school"
	Shop Building = "shop"
	TownCentre Building = "town_centre"
)

type Education string
const (
	Chef Education = "chef"
	Guard Education = "guard"
	Publicist Education = "publicist"
	Salesmouse Education = "salesmouse"
	Thief Education = "thief"
)
