package repository

import (
	"booking/internal/models"
	"booking/internal/storage"
	"context"
)

type HotelsRepository struct {
	*BaseRepository[models.HotelsModel]
}

func NewHotelsRepository(db *storage.Storage) *HotelsRepository {
	return &HotelsRepository{
		BaseRepository: NewBaseRepository[models.HotelsModel](db.DB),
	}
}

func (h *HotelsRepository) CreateHotel(ctx context.Context, hotels *models.HotelsModel) error {
	err := h.Add(ctx, hotels)
	if err != nil {
		return err
	}
	return nil
}
