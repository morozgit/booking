package service

import (
	"booking/internal/dto"
	"booking/internal/models"
	"booking/internal/repository"
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserNotFound = errors.New("user not found")
var ErrUnauthorized = errors.New("unauthorized")

type UserService struct {
	repo *repository.UserRepository
	log  *slog.Logger
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *UserService) VerifyPassword(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func (s *UserService) CreateAccessToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *UserService) GetUserIDFromToken(r *http.Request) (uint, error) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return 0, ErrUnauthorized
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return 0, ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrUnauthorized
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, ErrUnauthorized
	}

	return uint(userID), nil
}

func (s *UserService) RegisterUser(ctx context.Context, req dto.UserRequestAdd) error {

	hashedPassword, err := s.HashPassword(req.Password)
	if err != nil {
		s.log.Error("Failed to hash password", "error", err)
		return err
	}

	user := &models.UserModel{
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) Login(req dto.UserRequestAdd) (*models.UserModel, error) {
	user, err := s.repo.GetUser(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !s.VerifyPassword(req.Password, user.HashedPassword) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) GetByID(userID uint) (dto.User, error) {
	user, err := s.repo.GetOneOrNone(map[string]any{
		"id": userID,
	})
	if err != nil {
		return dto.User{}, err
	}

	if user == nil {
		return dto.User{}, ErrUserNotFound
	}

	return dto.User{
		Id:    user.ID,
		Email: user.Email,
	}, nil
}
