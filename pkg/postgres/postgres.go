package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/rgurov/bookkeeper/pkg/logging"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

const (
	delay    = time.Second * 5
	attempts = 5
)

type PostgresClient interface {
	Begin(context.Context) (pgx.Tx, error)
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) error
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

func Connect(ctx context.Context, cfg *Config) (pool *pgxpool.Pool, err error) {
	logger := logging.LoggerWithContext(ctx)
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password,
		cfg.Host, cfg.Port, cfg.Database,
	)

	err = doWithAttempts(
		func() error {
			logger.Info("trying to connect to postgres")
			ctx, cancel := context.WithTimeout(ctx, delay)
			defer cancel()

			pgCfg, err := pgxpool.ParseConfig(dsn)
			if err != nil {
				logger.Error("error parse postgres dsn")
				return err
			}
			pool, err = pgxpool.ConnectConfig(ctx, pgCfg)
			if err != nil {
				logger.Error("error connect postgres")
				return err
			}
			return nil
		},
		attempts,
		delay,
	)
	return pool, err
}

func doWithAttempts(f func() error, n int, delay time.Duration) (err error) {
	for i := 0; i < n; i++ {
		err = f()
		if err == nil {
			return nil
		}
		time.Sleep(delay)
	}
	return err
}
