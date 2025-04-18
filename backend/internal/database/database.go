package database

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UUID(v uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: v,
		Valid: true,
	}
}

func genRand() (string, error) {
	randInt, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		return "", err
	}
	baseline := big.NewInt(100000)
	if randInt.Cmp(baseline) == -1 {
		return randInt.Add(randInt, baseline).String(), nil
	}
	return randInt.String(), nil
}

type PostgresDriver struct {
	pool *pgxpool.Pool
}

func (pd *PostgresDriver) Close() {
	pd.pool.Close()
}

func CreatePostgresDriver(connectionURL string) (*PostgresDriver, error) {
	pool, err := pgxpool.New(context.Background(), connectionURL)
	if err != nil {
		return nil, err
	}

	return &PostgresDriver{
		pool: pool,
	}, nil
}
