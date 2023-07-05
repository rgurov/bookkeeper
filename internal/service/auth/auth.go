package auth

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/pkg/errors"

	"github.com/rgurov/bookkeeper/internal/domain"
	repoErrors "github.com/rgurov/bookkeeper/internal/repository/errors"
	serviceErrors "github.com/rgurov/bookkeeper/internal/service/errors"
)

type UserRepository interface {
	Create(ctx context.Context, login, password string) (int, error)
	One(ctx context.Context, id int) (*domain.User, error)
	Exists(ctx context.Context, login, password string) (int, error)
}

type AuthService struct {
	repo UserRepository
}

func NewAuthService(repo UserRepository) *AuthService {
	return &AuthService{repo}
}

func hash(input string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
}

func (s *AuthService) Register(ctx context.Context, login, password string) (int, error) {
	password = hash(password)
	id, err := s.repo.Create(ctx, login, password)
	if err != nil {
		if errors.Is(err, repoErrors.ErrUniqueConstaint) {
			return 0, serviceErrors.ErrDataBusy
		}
		return 0, err
	}
	return id, nil
}

func (s *AuthService) Login(ctx context.Context, login, password string) (int, error) {
	password = hash(password)
	id, err := s.repo.Exists(ctx, login, password)
	if err != nil {
		if errors.Is(err, repoErrors.ErrNotFound) {
			return 0, serviceErrors.ErrNotFound
		}
		return 0, err
	}
	return id, nil
}
