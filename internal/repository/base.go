package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{
		db: db,
	}
}

func (r *BaseRepository[T]) Add(ctx context.Context, model *T) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *BaseRepository[T]) GetOneOrNone(ctx context.Context, filters map[string]any) (*T, error) {
	var model T

	query := r.db.WithContext(ctx)

	for field, value := range filters {
		query = query.Where(field+" = ?", value)
	}

	err := query.First(&model).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &model, nil
}
