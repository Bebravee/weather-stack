package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/chup1x/weather-stack/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClothesRepository struct {
	db *gorm.DB
}

func NewClothesRepository(db *gorm.DB) *ClothesRepository {
	return &ClothesRepository{db: db}
}

func (r *ClothesRepository) CreateClothes(ctx context.Context, new *domain.WeatherClothesEntity) error {
	// Upsert по ключу city_id: переписываем path, если запись уже есть.
	return r.db.WithContext(ctx).
		Table("clothes").
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "city_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"path", "created_at"}),
		}).
		Create(new).Error
}

func (r *ClothesRepository) GetClothesByCode(ctx context.Context, code string) (*domain.WeatherClothesEntity, error) {
	clothes := &domain.WeatherClothesEntity{}
	startOfDay := time.Now().Truncate(24 * time.Hour)

	if err := r.db.WithContext(ctx).
		Table("clothes").
		Where("city_id = ? AND created_at >= ?", code, startOfDay).
		Order("created_at DESC").
		First(clothes).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrClothesNotFound
		}
		return nil, err
	}

	return clothes, nil
}
