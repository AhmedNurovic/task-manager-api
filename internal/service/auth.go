package service

import (
	"context"
	"errors"

	"time"

	"github.com/ahmednurovic/task-manager-api/internal/model"
	"github.com/ahmednurovic/task-manager-api/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthServicer interface {
	Register(ctx context.Context, email, password string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) AuthServicer {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*model.User, error) {
	existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := CreateToken(user.ID, s.jwtSecret)
	if err != nil {
		return "", ErrTokenGeneration
	}

	return token, nil
}

func CreateToken(userID uint, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})
	return token.SignedString([]byte(secret))
}
