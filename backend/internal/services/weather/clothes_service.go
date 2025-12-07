package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chup1x/weather-stack/internal/domain"
)

type clothesStorage interface {
	CreateClothes(context.Context, *domain.WeatherClothesEntity) error
	GetClothesByCode(context.Context, string) (*domain.WeatherClothesEntity, error)
}

type userByTelegramProvider interface {
	SelectByTelegramID(context.Context, int) (*domain.UserEntity, error)
}

type weatherByCityProvider interface {
	GetWeatherWithCache(context.Context, string) (*domain.WeatherEntity, error)
}

type ClothesService struct {
	clothesRepo clothesStorage
	users       userByTelegramProvider
	weather     weatherByCityProvider
	storageDir  string
	now         func() time.Time
	llm         llmClient
}

func NewClothesService(repo clothesStorage, users userByTelegramProvider, weather weatherByCityProvider, llm llmClient) *ClothesService {
	return &ClothesService{
		clothesRepo: repo,
		users:       users,
		weather:     weather,
		storageDir:  ".",
		now:         time.Now,
		llm:         llm,
	}
}

func (s *ClothesService) GetClothes(ctx context.Context, telegramID int) (map[string]any, error) {
	user, err := s.users.SelectByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("select user by telegram id: %w", err)
	}

	weather, err := s.weather.GetWeatherWithCache(ctx, user.CityW)
	if err != nil {
		return nil, fmt.Errorf("get weather for clothes: %w", err)
	}

	code := buildClothesCode(user, weather)

	existing, err := s.clothesRepo.GetClothesByCode(ctx, code)
	if err == nil && existing.Path != "" {
		body, readErr := os.ReadFile(existing.Path)
		if readErr == nil {
			out := map[string]any{}
			if jsonErr := json.Unmarshal(body, &out); jsonErr == nil {
				return out, nil
			}
		}
	}

	recommendation := map[string]any{}
	if s.llm != nil && s.llm.Enabled() {
		if text, llmErr := s.llm.RecommendClothes(ctx, s.llmInputFrom(user, weather)); llmErr == nil && text != "" {
			recommendation = s.buildLLMRecommendation(code, user, weather, text)
		}
	}
	if len(recommendation) == 0 {
		recommendation = s.stubRecommendation(code, user, weather)
	}

	fileName := fmt.Sprintf("clothes_%s.json", code)
	if s.storageDir != "" {
		fileName = filepath.Join(s.storageDir, fileName)
	}

	raw, _ := json.Marshal(recommendation)
	if writeErr := os.WriteFile(fileName, raw, 0o644); writeErr != nil {
		return nil, fmt.Errorf("save clothes stub: %w", writeErr)
	}

	_ = s.clothesRepo.CreateClothes(ctx, &domain.WeatherClothesEntity{
		Code: code,
		Path: fileName,
	})

	return recommendation, nil
}

func buildClothesCode(user *domain.UserEntity, weather *domain.WeatherEntity) string {
	// user prefs: comf (t-shirt), tol (hoodie), puh (jacket)
	p1 := bucket(user.TComfort, 30, 2, 15)
	p2 := bucket(user.TTol, 30, 4, 15)
	p3 := bucket(user.TPuh, 30, 4, 15)
	weatherCode := bucket(int(weather.Temperature), 60, 3, 20)
	return fmt.Sprintf("%02d%02d%02d%02d", p1, p2, p3, weatherCode)
}

func bucket(value int, start int, step int, maxBuckets int) int {
	// Higher temps map to lower bucket numbers; clamp both ends.
	diff := start - value
	if diff < 0 {
		return 1
	}
	idx := diff/step + 1
	if idx < 1 {
		idx = 1
	}
	if idx > maxBuckets {
		idx = maxBuckets
	}
	return idx
}

func (s *ClothesService) stubRecommendation(code string, user *domain.UserEntity, weather *domain.WeatherEntity) map[string]any {
	return map[string]any{
		"stub":         true,
		"code":         code,
		"message":      "LLM запросы отключены или не настроены. Это безопасная заглушка.",
		"how_to_test":  "Включите LLM (LLM_ENABLED=true и остальные переменные), повторите запрос и сравните ответ. Заглушку можно удалить после проверки.",
		"user_city":    user.CityW,
		"user_temps":   map[string]int{"comf": user.TComfort, "tol": user.TTol, "puh": user.TPuh},
		"weather_used": weather,
		"created_at":   s.now().Format(time.RFC3339),
	}
}

func (s *ClothesService) buildLLMRecommendation(code string, user *domain.UserEntity, weather *domain.WeatherEntity, text string) map[string]any {
	return map[string]any{
		"stub":         false,
		"code":         code,
		"message":      text,
		"user_city":    user.CityW,
		"user_temps":   map[string]int{"comf": user.TComfort, "tol": user.TTol, "puh": user.TPuh},
		"weather_used": weather,
		"created_at":   s.now().Format(time.RFC3339),
	}
}

func (s *ClothesService) llmInputFrom(user *domain.UserEntity, weather *domain.WeatherEntity) llmClothesInput {
	return llmClothesInput{
		City:        user.CityW,
		Temperature: weather.Temperature,
		Description: weather.Description,
		Humidity:    weather.Humidity,
		WindSpeed:   weather.WindSpeed,
		UserTemps:   map[string]int{"comf": user.TComfort, "tol": user.TTol, "puh": user.TPuh},
	}
}
