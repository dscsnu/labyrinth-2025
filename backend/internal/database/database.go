package database

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgconn"
)

// Attempt to generate a random 6 digit number and return a string representation
func genRand() (string, error) {

	randInt, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {

		return "", err
	}

	baseline := big.NewInt(100000)

	if randInt.Cmp(baseline) == -1 {

		return randInt.Add(baseline).String(), nil

	}

	return randInt.String(), nil

}

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

// teamCreator is the team member who creates the team
func (pd *PostgresDriver) CreateTeam(ctx context.Context, teamName string, teamCreator string) (created bool, err error) {

	// Attempt 5 times in total to ensure consistency and account for
	// any instances of duplicate teamIds or failed random number generation

	var pgErr *pgconn.PgError

	for range 5 {

		teamId, err := genRand()
		if err != nil {

			continue

		}

		_, err = pd.conn.Exec(ctx, "INSERT INTO Team(id, name) VALUES ($1, $2)", teamId, teamCreator)

		if err == nil {

			created = true
			break

		}
		// "23505" is the postgres uniqueness violation code
		if ok := errors.As(err, &pgErr); ok && pgErr.Code != "23505" {

			return created, err

		}

	}
	return created, nil

}
