package repository

import (
	"github.com/rgurov/bookkeeper/internal/repository/user"
	"github.com/rgurov/bookkeeper/pkg/postgres"
)

type repository struct {
	User *user.UserRepository
}

func NewRepository(client postgres.PostgresClient) *repository {
	return &repository{
		user.NewUserRespository(client),
	}
}
