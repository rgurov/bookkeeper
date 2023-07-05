package user

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/rgurov/bookkeeper/internal/domain"
	"github.com/rgurov/bookkeeper/internal/repository/errors"
	"github.com/rgurov/bookkeeper/pkg/postgres"
	"github.com/rgurov/pgerrors"
)

type UserRepository struct {
	client postgres.PostgresClient
}

func NewUserRespository(client postgres.PostgresClient) *UserRepository {
	return &UserRepository{client}
}

func (r *UserRepository) Create(ctx context.Context, login, password string) (int, error) {
	sql := `
	insert into users (login, password)
	values ($1, $2)
	returning id
	`
	var id int
	err := r.client.QueryRow(
		ctx,
		sql,
		login,
		password,
	).Scan(
		&id,
	)
	if err != nil {
		if pgerrors.Is(err, pgerrors.UniqueViolation) {
			return 0, errors.ErrUniqueConstaint
		}
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) One(ctx context.Context, id int) (*domain.User, error) {
	sql := `
	select
		id,
		login
	from users
	where
		id = $1
	`
	var user domain.User
	err := r.client.QueryRow(
		ctx,
		sql,
		id,
	).Scan(
		&user.ID,
		&user.Login,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Exists(ctx context.Context, login, password string) (int, error) {
	sql := `
	select
		id
	from users
	where
		login = $1 and
		password = $2
	`
	var id int
	err := r.client.QueryRow(
		ctx,
		sql,
		login,
		password,
	).Scan(
		&id,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, errors.ErrNotFound
		}
		return 0, err
	}
	return id, nil
}
