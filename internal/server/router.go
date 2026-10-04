package server

import (
	"booking/internal/handlers"
	mwLogger "booking/internal/middleware/logger"
	"net/http"

	_ "booking/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/exp/slog"
)

type Handlers struct {
	Auth *handlers.UserHandler
	//Booking  *booking.Handler
	//Room     *room.Handler
	Hotel *handlers.HotelsHandler
	//Facility *facility.Handler
	//Image    *image.Handler
}

func NewRouter(log *slog.Logger, h Handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},
		AllowCredentials: true,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(mwLogger.New(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Auth.RegisterUser)
		r.Post("/login", h.Auth.Login)
		r.Get("/me", h.Auth.GetMe)
		r.Post("/logout", h.Auth.Logout)
	})

	r.Route("/hotels", func(r chi.Router) {
		r.Post("/", h.Hotel.CreteHotel)

	})

	r.Get("/swagger/*", httpSwagger.Handler())
	return r
}
