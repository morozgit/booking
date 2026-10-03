package repository

import (
	"booking/internal/models"
	"booking/internal/storage"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type UserRepository struct {
	*BaseRepository[models.UserModel]
}

func NewUserRepository(db *storage.Storage) *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository[models.UserModel](db.DB),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.UserModel) error {
	err := r.Add(ctx, user)
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return ErrUserAlreadyExists
		}
	}

	return err
}

func (r *UserRepository) GetUser(email string) (*models.UserModel, error) {
	var user models.UserModel

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
