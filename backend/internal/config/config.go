package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   serverConfig
	Postgres postgresConfig
	News     newsServiceConfig
	Weather  weatherServiceConfig
	LLM      LLMConfig
}

type serverConfig struct {
	Port string `env:"SERVER_PORT" envDefault:"80"`
}

type postgresConfig struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
	Name     string `env:"POSTGRES_NAME" envDefault:"postgres"`
	User     string `env:"POSTGRES_USER" envDefault:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	SSLMode  string `env:"POSTGRES_SSLMODE" envDefault:"disable"`
}

type newsServiceConfig struct {
	Host   string `env:"NEW_SERVICE_HOST" envDefault:"newsapi.org/v2/everything"`
	APIKey string `env:"NEW_SERVICE_API_KEY" envDefault:"0fac40f7dcd34967af176019e1c6a526"`
}

type weatherServiceConfig struct {
	Host   string `env:"WEATHER_SERVICE_HOST" envDefault:"api.openweathermap.org/data/2.5/weather"`
	APIKey string `env:"WEATHER_SERVICE_API_KEY" envDefault:""`
	Lang   string `env:"WEATHER_LANG" envDefault:"ru"`
	Units  string `env:"WEATHER_UNITS" envDefault:"metric"`
}

type LLMConfig struct {
	Enabled     bool    `env:"LLM_ENABLED" envDefault:"false"`
	URL         string  `env:"LLM_URL"`
	APIKey      string  `env:"LLM_API_KEY"`
	Model       string  `env:"LLM_MODEL"`
	Temperature float64 `env:"LLM_TEMPERATURE" envDefault:"0.7"`
	MaxTokens   int     `env:"LLM_MAX_TOKENS" envDefault:"400"`
	TimeoutSec  int     `env:"LLM_TIMEOUT_SEC" envDefault:"10"`
}

func GetConfig() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{}
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("parse .env file: %w", err)
	}

	return config, nil
}
