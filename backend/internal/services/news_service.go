package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"io"
	"encoding/json"
	"os"
	"github.com/chup1x/weather-stack/internal/domain"
)

type newsProvider interface {
	CreateNewsRequest(context.Context) (*domain.NewsEntity, error)
	GetNewsByCityID(context.Context, string) (*domain.NewsEntity, error)
}

type newsStorage interface {
	newsProvider
}

type NewsService struct {
	repo newsStorage
}

func NewNewsService(repo newsStorage) *NewsService {
	return &NewsService{repo: repo}
}

func (s *NewsService) GetNews(ctx context.Context, cityID string) (*domain.NewsEntity, error) {
	news, err := s.repo.GetNewsByCityID(ctx, cityID)
	if err != nil {
	 	fmt.Errorf("to select a weather by city: %w", err)
	}
	if errors.Is(err, domain.NewsNotFound) {
	

		baseURL := "https://newsapi.org/v2/everything"
		params := url.Values{}
		params.Add("q", cityID)
		params.Add("from", "2025-11-24")
		params.Add("sortBy", "publishedAt")
		params.Add("language", "ru")
		params.Add("apiKey", "0fac40f7dcd34967af176019e1c6a526")

		fullURL := baseURL + "?" + params.Encode()

		resp, err := http.Get(fullURL)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		filename := fmt.Sprintf("temp_news_%s.json", cityID)
		err = os.WriteFile(filename, body, 0644)
		if err != nil {
			log.Fatal("Error writing file:", err)
		}
		news_en, _ := json.Marshal(body)
		news := &domain.NewsEntity{
			CityID : cityID,
			Path: filename,
		}
		if err := s.repo.CreateNewsRequest(ctx, news); err != nil {
			return news, fmt.Errorf("to create a weatcher request: %w", err)
		}
		return news, nil
	}
	return news, nil
}
