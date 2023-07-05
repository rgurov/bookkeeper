package service

import "github.com/rgurov/bookkeeper/internal/service/auth"

type service struct {
	Auth *auth.AuthService
}

func NewService(
	userRepo auth.UserRepository,
) *service {
	return &service{
		Auth: auth.NewAuthService(userRepo),
	}
}
