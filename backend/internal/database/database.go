package database

import (
	"context"
	"crypto/rand"
	"errors"
	"labyrinth/internal/types"
	"math/big"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func UUID(v uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: v,
		Valid: true,
	}
}

// Attempt to generate a random 6 digit number and return a string representation
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

// teamCreatorId is the team member who creates the team
func (pd *PostgresDriver) CreateTeam(ctx context.Context, teamName string, teamCreatorId uuid.UUID) (bool, error) {

	// Attempt 5 times in total to ensure consistency and account for
	// any instances of duplicate teamIds or failed random number generation

	var pgErr *pgconn.PgError

	var teamCreated bool
	var teamCreatedId string

	for range 5 {

		teamId, err := genRand()
		if err != nil {

			continue

		}

		_, err = pd.conn.Exec(ctx, "INSERT INTO team(id, name) VALUES ($1, $2)", teamId, teamName)

		if err == nil {

			teamCreated = true
			teamCreatedId = teamId
			break

		}
		// "23505" is the postgres uniqueness violation code
		if ok := errors.As(err, &pgErr); ok && pgErr.Code != "23505" {

			return false, err

		}

	}

	if !teamCreated {

		return false, errors.New("Team could not be created due to internal error")

	}

	pgxUuid := UUID(teamCreatorId)

	if _, err := pd.conn.Exec(ctx, "INSERT INTO teammember(team_id, user_id) VALUES ($1,$2)", teamCreatedId, pgxUuid); err != nil {

		return false, err

	}

	return true, nil

}

func (pd *PostgresDriver) GetUser(ctx context.Context, userEmail string) (types.UserProfile, error) {

	userProfile := types.UserProfile{}

	row := pd.conn.QueryRow(ctx, "SELECT id, name, email, created_at from userprofile WHERE email=$1", userEmail)
	if err := row.Scan(&userProfile.ID, &userProfile.Name, &userProfile.Email, &userProfile.CreatedAt); err != nil {

		return userProfile, err

	}

	return userProfile, nil
}
