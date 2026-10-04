package service

import (
	"booking/internal/dto"
	"booking/internal/models"
	"booking/internal/repository"
	"context"
	"errors"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserNotFound = errors.New("user not found")
var ErrUnauthorized = errors.New("unauthorized")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
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
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return "", errors.New("JWT_SECRET_KEY is not set")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *UserService) GetUserIDFromToken(r *http.Request) (uint, error) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return 0, ErrUnauthorized
	}
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return 0, ErrUnauthorized
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return 0, ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrUnauthorized
	}

	userID, ok := claims["user_id"].(float64)
	if !ok || userID <= 0 || userID >= float64(uint64(1)<<63) || math.Trunc(userID) != userID {
		return 0, ErrUnauthorized
	}

	return uint(userID), nil
}

func (s *UserService) RegisterUser(ctx context.Context, req dto.UserRequestAdd) error {

	hashedPassword, err := s.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &models.UserModel{
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) Login(ctx context.Context, req dto.UserRequestAdd) (*models.UserModel, error) {
	user, err := s.repo.GetUser(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if user == nil {
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
