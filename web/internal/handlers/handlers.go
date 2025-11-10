package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Provider struct {
	Host      string
	validator *validator.Validate
	Client    *http.Client
}

func New(host string) *Provider {
	return &Provider{
		Host:      host,
		validator: validator.New(),
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *Provider) RegisterUserHandler(c *fiber.Ctx) error {
	var req RegisterUserRequest

	if err := c.BodyParser(&req); err != nil {
		fmt.Errorf("bla", err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if err := p.validator.Struct(req); err != nil {
		fmt.Errorf("blab", err)
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	rawBody, err := json.Marshal(req)
	if err != nil {
		fmt.Errorf("blabl", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	body, statusCode, err := p.sendRequest(
		c.UserContext(),
		http.MethodPost,
		fmt.Sprintf("%s/auth/register", p.Host),
		bytes.NewReader(rawBody),
	)
	if err != nil {
		return c.SendStatus(statusCode)
	}

	res := RegisterUserResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Errorf("blabla", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(res)
}

func (p *Provider) GetUserHandler(c *fiber.Ctx) error {
	var req GetUserRequest
	if err := c.ParamsParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if err := p.validator.Struct(req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	rawBody, err := json.Marshal(map[string]any{"user_id": req.ID})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	body, statusCode, err := p.sendRequest(
		c.UserContext(),
		http.MethodGet,
		fmt.Sprintf("%s/profile/by-id", p.Host),
		bytes.NewReader(rawBody),
	)
	if err != nil {
		return c.SendStatus(statusCode)
	}

	res := GetUserResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(res)
}

func (p *Provider) sendRequest(ctx context.Context, method string, url string, body io.Reader) ([]byte, int, error) {
	req, _ := http.NewRequestWithContext(ctx, method, url, body)
	req.Header.Set("Content-Type", "application/json")

	res, err := p.Client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("to send a request: %w", err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("to read a body: %w", err)
	}

	return data, res.StatusCode, nil
}
