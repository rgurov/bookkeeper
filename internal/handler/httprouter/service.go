package httprouter

import "context"

type AuthService interface {
	Register(ctx context.Context, login, password string) (int, error)
	Login(ctx context.Context, login, password string) (int, error)
}
