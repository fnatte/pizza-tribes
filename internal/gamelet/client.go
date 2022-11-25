package gamelet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/fnatte/pizza-tribes/internal/game/protojson"
	"github.com/fnatte/pizza-tribes/internal/game/models"
)

type GameletClient struct {
	host       string
	httpClient *http.Client
}

func NewGameletClient(host string) *GameletClient {
	return &GameletClient{
		host:       host,
		httpClient: &http.Client{},
	}
}

func (gc *GameletClient) JoinGame(userId, username string, items []string) error {
	data, err := json.Marshal(map[string]interface{}{
		"user_id":  userId,
		"username": username,
		"items":    items,
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/join", gc.host)

	resp, err := gc.httpClient.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Got status code: %d", resp.StatusCode)
	}

	return nil
}

func (gc *GameletClient) GetLeaderboard(skip int) (*models.Leaderboard, error) {
	url := fmt.Sprintf("http://%s/leaderboard?skip=%d", gc.host, skip)

	resp, err := gc.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Got status code: %d", resp.StatusCode)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read body")
	}

	lb := &models.Leaderboard{}
	if err = protojson.Unmarshal(buf, lb); err != nil {
		return nil, errors.New("failed to parse body")
	}

	return lb, nil
}
