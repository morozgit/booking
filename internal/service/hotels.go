package service

import (
	"booking/internal/dto"
	"booking/internal/models"
	"booking/internal/repository"
	"context"
)

type HotelsService struct {
	repo *repository.HotelsRepository
}

func NewHotelsService(repo *repository.HotelsRepository) *HotelsService {
	return &HotelsService{repo: repo}
}

func (s *HotelsService) AddHotel(ctx context.Context, req dto.HotelAdd) error {
	hotel := &models.HotelsModel{
		Title:    req.Title,
		Location: req.Location,
	}
	return s.repo.Add(ctx, hotel)
}
