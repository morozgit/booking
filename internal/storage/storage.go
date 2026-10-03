package storage

import (
	"booking/internal/config"
	"fmt"

	"github.com/WqyJh/go-fstring"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	DB *gorm.DB
}

func New(cfg *config.Config) (*Storage, error) {
	const op = "storage.postgres.New"
	template := "host={host} user={user} password={password} dbname={dbname} port={port}"
	values := map[string]any{
		"host":     cfg.Host,
		"user":     cfg.User,
		"password": cfg.Password,
		"dbname":   cfg.DBName,
		"port":     cfg.Port,
	}
	dsn, err := fstring.Format(template, values)
	if err != nil {
		return nil, fmt.Errorf("%s - failed to build DSN: %w", op, err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s - failed to connect db: %w", op, err)
	}

	return &Storage{DB: db}, nil
}
