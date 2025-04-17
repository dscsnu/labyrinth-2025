package database

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"
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
	pool  *pgxpool.Pool
	cache *cache.Cache
}

func (pd *PostgresDriver) Close() {
	pd.pool.Close()
}

func CreatePostgresDriver(connectionURL string) (*PostgresDriver, error) {
	pool, err := pgxpool.New(context.Background(), connectionURL)
	if err != nil {
		return nil, err
	}

	c := cache.New(5*time.Minute, 10*time.Minute)

	return &PostgresDriver{
		pool:  pool,
		cache: c,
	}, nil
}

func (pd *PostgresDriver) SetCacheTTL(defaultExpiration, cleanupInterval time.Duration) {
	pd.cache = cache.New(defaultExpiration, cleanupInterval)
}

func (pd *PostgresDriver) CacheSet(key string, value interface{}, expiration time.Duration) {
	pd.cache.Set(key, value, expiration)
}

func (pd *PostgresDriver) CacheGet(key string) (interface{}, bool) {
	return pd.cache.Get(key)
}

func (pd *PostgresDriver) CacheDelete(key string) {
	pd.cache.Delete(key)
}

func (pd *PostgresDriver) CacheFlush() {
	pd.cache.Flush()
}
