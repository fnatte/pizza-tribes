package main

import (
	"net/http"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/fnatte/pizza-tribes/internal/redis"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type TimeseriesService struct {
	r    redis.RedisClient
	auth *internal.AuthService
}

func (s *TimeseriesService) Handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		err := s.auth.Authorize(r)
		if err != nil {
			log.Error().Err(err).Msg("Failed to authorize")
			w.WriteHeader(403)
			return
		}
		userId, ok := r.Context().Value("userId").(string)
		if !ok {
			log.Warn().Msg("Failed to get account id")
			w.WriteHeader(500)
			return
		}

		tsPizzas, err := internal.FetchPizzasTimeseries(r.Context(), s.r, userId)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch pizzas timeseries")
			w.WriteHeader(500)
			return
		}
		tsCoins, err := internal.FetchCoinsTimeseries(r.Context(), s.r, userId)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch pizzas timeseries")
			w.WriteHeader(500)
			return
		}

		ts := mergeTimeseries(tsPizzas, tsCoins)
		d := &models.TimeseriesData{
			DataPoints: ts,
		}
		b, err := protojson.Marshal(d)
		if err != nil {
			w.WriteHeader(500)
			log.Error().Err(err).Msg("Failed to marshal timeseries data")
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(b)
	})

	return r
}

// Merge timeseries data and convert from Redis format to message format.
func mergeTimeseries(tsPizzas []*redis.TimeseriesDataPoint, tsCoins []*redis.TimeseriesDataPoint) []*models.DataPoint {
	res := []*models.DataPoint{}
	pi := 0
	ci := 0
	for pi < len(tsPizzas) && ci < len(tsCoins) {
		p := tsPizzas[pi]
		c := tsCoins[ci]

		if p.Timestamp == c.Timestamp {
			res = append(res, &models.DataPoint{
				Timestamp: p.Timestamp,
				Coins:     c.Value,
				Pizzas:    p.Value,
			})
		} else {
			dpp := &models.DataPoint{
				Timestamp: p.Timestamp,
				Pizzas:    p.Value,
			}
			dpc := &models.DataPoint{
				Timestamp: p.Timestamp,
				Coins:     c.Value,
			}
			if dpp.Timestamp < dpc.Timestamp {
				res = append(res, dpp, dpc)
			} else {
				res = append(res, dpc, dpp)
			}
		}

		ci++
		pi++
	}
	for pi < len(tsPizzas) {
		res = append(res, &models.DataPoint{
			Timestamp: tsPizzas[pi].Timestamp,
			Pizzas:    tsPizzas[pi].Value,
		})
		pi++
	}
	for ci < len(tsCoins) {
		res = append(res, &models.DataPoint{
			Timestamp: tsCoins[ci].Timestamp,
			Coins:     tsCoins[ci].Value,
		})
		ci++
	}

	return res
}
