package storage

import (
	"context"
	"log/slog"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)


type Postgres struct{
	Pool *pgxpool.Pool
	logger *slog.Logger
}

//postgresql://neondb_owner:npg_CETJIpY68ZvW@ep-frosty-math-b30xk6i1-pooler.c-4.ap-southeast-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require

func NewPostgres(ctx context.Context, logger *slog.Logger, databaseURL string) (*Postgres, error)  {
	pool, err:= pgxpool.New(ctx, databaseURL)
	if err!=nil{
		logger.Error("unable to create a new pool", "error", err)
		return nil, fmt.Errorf("unable to create pool: %w", err)
	}

	if err:=pool.Ping(ctx); err!=nil {
		logger.Error("unable to ping db", "error", err)
		pool.Close()
		return nil, fmt.Errorf("unable to ping pool:%w", err)
	}

	return &Postgres{
		Pool: pool,
		logger: logger,
	}, nil

}