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

type newsCreater interface {
	CreateNewsRequest(context.Context, *domain.NewsEntity) error
}

type newsProvider interface {
	GetNewsByCityID(context.Context, string) (*domain.NewsEntity, error)
}

type newsStorage interface {
	newsCreater
	newsProvider
}

type NewsService struct {
	repo       newsStorage
	newsServ   *newsClient
	storageDir string
	now        func() time.Time
}

func NewNewsService(repo newsStorage, news *newsClient) *NewsService {
	return &NewsService{
		repo:       repo,
		newsServ:   news,
		storageDir: ".",
		now:        time.Now,
	}
}

func (s *NewsService) GetNews(ctx context.Context, cityID string) (map[string]any, error) {
	// Сначала пытаемся отдать кеш (файл по пути из таблицы news).
	existing, err := s.repo.GetNewsByCityID(ctx, cityID)
	if err == nil && existing.Path != "" {
		body, readErr := os.ReadFile(existing.Path)
		if readErr == nil {
			out := map[string]any{}
			if jsonErr := json.Unmarshal(body, &out); jsonErr == nil {
				return out, nil
			}
		}
	}

	// Запрашиваем свежие новости у внешнего API.
	news, err := s.newsServ.GetNewsByCityID(ctx, cityID)
	if err != nil {
		return nil, fmt.Errorf("to get a news by city id - %s: %w", cityID, err)
	}

	// Сохраняем ответ в файл и записываем путь в БД.
	fileName := fmt.Sprintf("news_%s_%s.json", cityID, s.now().Format("2006-01-02"))
	if s.storageDir != "" {
		fileName = filepath.Join(s.storageDir, fileName)
	}

	raw, _ := json.Marshal(news)
	if writeErr := os.WriteFile(fileName, raw, 0o644); writeErr != nil {
		return nil, fmt.Errorf("to save news to file: %w", writeErr)
	}

	_ = s.repo.CreateNewsRequest(ctx, &domain.NewsEntity{
		CityID:    cityID,
		Path:      fileName,
		CreatedAt: s.now(),
	})

	return news, nil
}
