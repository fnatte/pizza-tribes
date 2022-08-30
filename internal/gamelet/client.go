package gamelet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
		"items": items,
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
