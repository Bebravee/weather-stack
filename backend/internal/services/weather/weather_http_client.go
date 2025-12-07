package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/chup1x/weather-stack/internal/domain"
)

type weatherClient struct {
	host   string
	apiKey string
	lang   string
	units  string

	client *http.Client
}

func NewWeatherClient(host, apiKey, lang, units string) *weatherClient {
	return &weatherClient{
		host:   host,
		apiKey: apiKey,
		lang:   lang,
		units:  units,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *weatherClient) FetchWeather(ctx context.Context, city string) (*domain.WeatherEntity, error) {
	values := url.Values{}
	values.Add("appid", c.apiKey)
	values.Add("q", city)
	if c.units != "" {
		values.Add("units", c.units)
	}
	if c.lang != "" {
		values.Add("lang", c.lang)
	}

	addr := fmt.Sprintf("https://%s?%s", c.host, values.Encode())

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("weather request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather request failed: status %d", res.StatusCode)
	}

	body := struct {
		Main struct {
			Temp     float64 `json:"temp"`
			Feels    float64 `json:"feels_like"`
			Humidity float64 `json:"humidity"`
			Pressure float64 `json:"pressure"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
		Wind struct {
			Speed float64 `json:"speed"`
		} `json:"wind"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode weather body: %w", err)
	}

	description := ""
	if len(body.Weather) > 0 {
		description = body.Weather[0].Description
	}

	return &domain.WeatherEntity{
		CityID:      city,
		Temperature: body.Main.Temp,
		FeelsLike:   body.Main.Feels,
		Description: description,
		Humidity:    body.Main.Humidity,
		Pressure:    body.Main.Pressure,
		WindSpeed:   body.Wind.Speed,
		CreatedAt:   time.Now(),
	}, nil
}
