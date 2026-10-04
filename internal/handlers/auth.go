package handlers

import (
	"booking/internal/dto"
	"booking/internal/repository"
	"booking/internal/service"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

type UserHandler struct {
	service *service.UserService
	vld     *validator.Validate
	log     *slog.Logger
}

func NewUserHandler(service *service.UserService, vld *validator.Validate, log *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		vld:     vld,
		log:     log,
	}
}

// RegisterUser @Summary      Create User
// @Description  Add User with hashed password
// @Tags         auth
// @Produce      json
// @Success      200 {object} dto.StatusResponse
// @Param request body dto.UserRequestAdd true "Registration data"
// @Router       /auth/register [post]
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.UserRequestAdd

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode request", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.RegisterUser(ctx, req); err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		h.log.Error("failed to register user", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "OK",
	})
}

// Login @Summary      Login User
// @Description  Returns JWT token
// @Tags         auth
// @Produce      json
// @Success      200 {object} dto.AccessTokenResponse
// @Param request body dto.UserRequestAdd true "Registration data"
// @Router       /auth/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.UserRequestAdd

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode request", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(ctx, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		h.log.Error("failed to login", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	accessToken, err := h.service.CreateAccessToken(user.ID)
	if err != nil {
		h.log.Error("failed to create access token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   24 * 60 * 60,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}

// GetMe @Summary      Get User
// @Description  Returns User
// @Tags         auth
// @Produce      json
// @Success      200 {object} dto.User
// @Router       /auth/me [get]
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, err := h.service.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetByID(userID)
	if err != nil {
		h.log.Error("failed to get current user", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]any{
		"User": user,
	})
}

// Logout @Summary      Logout
// @Description  Logout
// @Tags         auth
// @Produce      json
// @Success      200 {object} dto.StatusResponse
// @Router       /auth/logout [post]
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"Status": "OK",
	})

}
