package app

import (
	"booking/internal/config"
	"booking/internal/handlers"
	"booking/internal/lib/logger/sl"
	"booking/internal/lib/validator"
	"booking/internal/repository"
	"booking/internal/server"
	"booking/internal/service"
	"booking/internal/storage"
	"context"
	"os"

	"golang.org/x/exp/slog"
)

type App struct {
	server *server.Server
	db     *storage.Storage
	//redis  *redis.Client
}

func New(cfg *config.Config, log *slog.Logger) (*App, error) {

	db, err := storage.New(cfg)
	if err != nil {
		log.Error("Failed to initialize storage:", sl.Err(err))
		os.Exit(1)
	}

	validate := validator.New()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService, validate, log)

	hotelRepository := repository.NewHotelsRepository(db)
	hotelService := service.NewHotelsService(hotelRepository)
	hotelHandler := handlers.NewHotelsHandler(hotelService)

	router := server.NewRouter(log, server.Handlers{
		Auth:  userHandler,
		Hotel: hotelHandler,
	})
	srv := server.New(cfg, router, log)

	return &App{
		server: srv,
		db:     db,
	}, nil
}

func (a *App) Run() error {
	return a.server.Run()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
