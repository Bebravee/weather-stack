package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type newsClient struct {
	host   string
	apiKey string

	client *http.Client
}

func NewNewsClient(host string, apiKey string) *newsClient {
	return &newsClient{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		host:   host,
		apiKey: apiKey,
	}
}

func (n *newsClient) GetNewsByCityID(ctx context.Context, cityID string) (map[string]any, error) {
	values := url.Values{}
	values.Add("apiKey", n.apiKey)
	values.Add("q", cityID)
	// Берем неделю назад, иначе за текущий день часто пусто
	values.Add("from", time.Now().AddDate(0, 0, -7).Format("2006-01-02"))
	values.Add("to", time.Now().Format("2006-01-02"))
	values.Add("sortBy", "publishedAt")
	values.Add("language", "ru")

	addr := fmt.Sprintf("https://%s?%s", n.host, values.Encode())

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
	res, err := n.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("to get a news response: %w", err)
	}
	defer res.Body.Close()

	body := make(map[string]any)
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("to decode json from response body: %w", err)
	}

	return body, nil
}
