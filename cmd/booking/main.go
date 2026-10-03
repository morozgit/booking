package main

import (
	"booking/internal/app"
	"booking/internal/config"
	"booking/internal/lib/logger/sl"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Booking API
// @version 1.0
// @description Booking service API
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.MustLoad()

	log := sl.SetupLogger(cfg.Env)

	_app, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to initialize application", sl.Err(err))
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		if err := _app.Run(); err != nil {
			log.Error("server error", sl.Err(err))
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := _app.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown error", sl.Err(err))
	}
}
