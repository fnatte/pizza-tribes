package internal

import (
	"math"
	"time"

	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/rs/zerolog/log"
)

func calculateQualityScore(gs *models.GameState) float64 {
	score := 1.0

	if gs.HasDiscovery(models.ResearchDiscovery_DURUM_WHEAT) {
		score = score + 0.05
	}
	if gs.HasDiscovery(models.ResearchDiscovery_DOUBLE_ZERO_FLOUR) {
		score = score + 0.05
	}
	if gs.HasDiscovery(models.ResearchDiscovery_SAN_MARZANO_TOMATOES) {
		score = score + 0.05
	}
	if gs.HasDiscovery(models.ResearchDiscovery_OCIMUM_BASILICUM) {
		score = score + 0.05
	}
	if gs.HasDiscovery(models.ResearchDiscovery_EXTRA_VIRGIN) {
		score = score + 0.05
	}
	if gs.HasDiscovery(models.ResearchDiscovery_MASONRY_OVEN) {
		score = score + 0.1
	}

	return score
}

func calculateMarketingScore(gs *models.GameState, e map[models.Education]int32) float64 {
	marketingBonus := 0.0

	if gs.HasDiscovery(models.ResearchDiscovery_WEBSITE) {
		marketingBonus += 0.1
	}
	if gs.HasDiscovery(models.ResearchDiscovery_MOBILE_APP) {
		marketingBonus += 0.1
	}

	publicists := float64(e[models.Education_PUBLICIST])

	return (1 + publicists) * (1 + marketingBonus)
}

const ep = -1.5 // Price exponant
const em = 1.2  // Marketing exponant
const er = 1.5  // Quality exponant

func CalculateDemandScore(gs *models.GameState) float64 {
	educations := CountTownPopulationEducations(gs)

	p := float64(gs.GetValidPizzaPrice())        // Price
	r := calculateQualityScore(gs)               // Quality
	m := calculateMarketingScore(gs, educations) // Marketing
	e := 1.0                                     // Economical index
	s := 1.0                                     // Seasonal index

	log.Debug().Float64("price", p).Float64("quality", r).Float64("marketing", m).Msg("Calculate demand score")

	score := math.Pow(p, ep) * math.Pow(r, er) * math.Pow(m, em) * e * s

	return score
}

func CalculateGlobalDemand(worldState *models.WorldState, userCount int64) float64 {
	start := time.Unix(worldState.StartTime, 0)
	dur := time.Now().Sub(start)
	days := dur.Hours() / 24

	d := days * float64(userCount)

	return d
}
