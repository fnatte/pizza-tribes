package main

import (
	"net/http"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

type TimeseriesService struct {
	r    internal.RedisClient
	auth *AuthService
}

func (s *TimeseriesService) Router() http.Handler {
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
		d := &internal.TimeseriesData{
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
func mergeTimeseries(tsPizzas []*internal.TimeseriesDataPoint, tsCoins []*internal.TimeseriesDataPoint) []*internal.DataPoint {
	res := []*internal.DataPoint{}
	pi := 0
	ci := 0
	for pi < len(tsPizzas) && ci < len(tsCoins) {
		p := tsPizzas[pi]
		c := tsCoins[ci]

		if p.Timestamp == c.Timestamp {
			res = append(res, &internal.DataPoint{
				Timestamp: p.Timestamp,
				Coins:     c.Value,
				Pizzas:    p.Value,
			})
		} else {
			dpp := &internal.DataPoint{
				Timestamp: p.Timestamp,
				Pizzas:    p.Value,
			}
			dpc := &internal.DataPoint{
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
		res = append(res, &internal.DataPoint{
			Timestamp: tsPizzas[pi].Timestamp,
			Pizzas:    tsPizzas[pi].Value,
		})
		pi++
	}
	for ci < len(tsCoins) {
		res = append(res, &internal.DataPoint{
			Timestamp: tsCoins[ci].Timestamp,
			Coins:     tsCoins[ci].Value,
		})
		ci++
	}

	return res
}
