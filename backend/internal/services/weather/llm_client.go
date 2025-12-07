package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chup1x/weather-stack/internal/config"
)

type llmClient interface {
	Enabled() bool
	RecommendClothes(ctx context.Context, in llmClothesInput) (string, error)
}

type llmClothesInput struct {
	City        string
	Temperature float64
	Description string
	Humidity    float64
	WindSpeed   float64
	UserTemps   map[string]int
}

func NewLLMClient(cfg config.LLMConfig) llmClient {
	if !cfg.Enabled || cfg.URL == "" || cfg.APIKey == "" || cfg.Model == "" {
		return disabledLLM{}
	}

	timeout := time.Duration(cfg.TimeoutSec) * time.Second
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &openAILLM{
		url:         cfg.URL,
		apiKey:      cfg.APIKey,
		model:       cfg.Model,
		temperature: cfg.Temperature,
		maxTokens:   cfg.MaxTokens,
		client:      &http.Client{Timeout: timeout},
	}
}

type disabledLLM struct{}

func (disabledLLM) Enabled() bool { return false }

func (disabledLLM) RecommendClothes(ctx context.Context, _ llmClothesInput) (string, error) {
	return "", errors.New("llm disabled")
}

type openAILLM struct {
	url         string
	apiKey      string
	model       string
	temperature float64
	maxTokens   int
	client      *http.Client
}

func (c *openAILLM) Enabled() bool { return true }

func (c *openAILLM) RecommendClothes(ctx context.Context, in llmClothesInput) (string, error) {
	systemPrompt := "Ты стилист, даешь краткие рекомендации по одежде, учитывая температуру, описание погоды, влажность и ветер. Отвечай по-русски."

	userLines := []string{
		fmt.Sprintf("Город: %s", in.City),
		fmt.Sprintf("Температура: %.1f°C, описание: %s", in.Temperature, in.Description),
		fmt.Sprintf("Влажность: %.0f%%, ветер: %.1f м/с", in.Humidity, in.WindSpeed),
	}
	if len(in.UserTemps) > 0 {
		userLines = append(userLines, fmt.Sprintf("Предпочтения пользователя (футболка/толстовка/пуховик): %v", in.UserTemps))
	}

	payload := map[string]any{
		"model":       c.model,
		"temperature": c.temperature,
		"max_tokens":  c.maxTokens,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": strings.Join(userLines, "\n")},
		},
	}

	raw, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("build llm request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	res, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do llm request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("llm request failed: status %d", res.StatusCode)
	}

	body := struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("decode llm response: %w", err)
	}

	if len(body.Choices) == 0 {
		return "", errors.New("llm response empty")
	}

	return strings.TrimSpace(body.Choices[0].Message.Content), nil
}
