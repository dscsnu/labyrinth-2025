package database

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type PostgresDriver struct {
	conn *pgx.Conn
}

func (pd *PostgresDriver) Close(ctx context.Context) error {

	return pd.conn.Close(ctx)

}

// Connect to postgres instance
func CreatePostgresDriver(connectionURL string) (*PostgresDriver, error) {

	conn, err := pgx.Connect(context.Background(), connectionURL)
	if err != nil {
		return &PostgresDriver{}, err
	}

	return &PostgresDriver{

		conn: conn,
	}, nil

}
