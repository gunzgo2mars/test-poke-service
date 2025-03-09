package auth

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/gunzgo2mars/test-poke-service/app/internal/core/model"
	"github.com/gunzgo2mars/test-poke-service/app/internal/repository/users"
	"github.com/gunzgo2mars/test-poke-service/app/pkg/middleware"
	"github.com/gunzgo2mars/test-poke-service/app/pkg/utils"
)

type IAuthService interface {
	ValidatingUser(ctx context.Context, req *model.AuthRequest) (string, error)
	RegisterUser(ctx context.Context, req *model.AuthRequest) error
}

type authService struct {
	// Dependencies.
	jwtMiddleware middleware.IJwt

	// Repositories.
	userRepo users.IUserRepository
}

func New(jwtMiddleware middleware.IJwt, userRepo users.IUserRepository) IAuthService {
	return &authService{
		jwtMiddleware: jwtMiddleware,
		userRepo:      userRepo,
	}
}

func (s *authService) ValidatingUser(ctx context.Context, req *model.AuthRequest) (string, error) {
	user, err := s.userRepo.GetUser(ctx, users.ByUsername(req.Username))
	if err != nil {
		return "", err
	}

	if err := utils.Compare([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", err
	}

	payload := jwt.MapClaims{}
	payload["uuid"] = user.UUID

	token, err := s.jwtMiddleware.GenerateToken(payload)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) RegisterUser(ctx context.Context, req *model.AuthRequest) error {
	now := time.Now()

	hash, err := utils.GenerateHash(req.Password)
	if err != nil {
		return err
	}

	schema := &model.UserSchema{
		UUID:     uuid.New().String(),
		Username: req.Username,
		Password: string(hash),
		CreateAt: now,
		UpdateAt: now,
	}

	return s.userRepo.CreateNewUser(ctx, schema)
}
