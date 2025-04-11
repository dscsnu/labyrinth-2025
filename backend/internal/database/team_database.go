package database

import (
	"context"
	"errors"
	"labyrinth/internal/types"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func (pd *PostgresDriver) CreateTeam(ctx context.Context, teamName string, teamCreatorId uuid.UUID) (string, error) {
	var pgErr *pgconn.PgError
	var teamCreatedId string

	for range 5 {
		teamId, err := genRand()
		if err != nil {
			continue
		}

		_, err = pd.pool.Exec(ctx, "INSERT INTO team(id, name) VALUES ($1, $2)", teamId, teamName)
		if err == nil {
			teamCreatedId = teamId
			break
		}
		if ok := errors.As(err, &pgErr); ok && pgErr.Code != "23505" {
			return "", err
		}
	}

	if teamCreatedId == "" {
		return "", errors.New("team could not be created due to internal error")
	}

	pgxUuid := UUID(teamCreatorId)
	if _, err := pd.pool.Exec(ctx, "INSERT INTO teammember(team_id, user_id) VALUES ($1,$2)", teamCreatedId, pgxUuid); err != nil {
		return "", err
	}

	return teamCreatedId, nil
}

func (pd *PostgresDriver) AddTeamMember(ctx context.Context, teamId string, userId uuid.UUID) error {
	var count int

	err := pd.pool.QueryRow(ctx, `SELECT COUNT(*) FROM teammember WHERE team_id=$1`, teamId).Scan(&count)

	if err != nil {
		return err
	}

	if count >= 4 {
		return errors.New("team is full")
	}

	if _, err := pd.pool.Exec(ctx, "INSERT INTO teammember(team_id, user_id) VALUES ($1,$2)", teamId, userId); err != nil {
		return err
	}

	return nil
}

func (pd *PostgresDriver) GetUser(ctx context.Context, userEmail string) (types.UserProfile, error) {
	userProfile := types.UserProfile{}
	row := pd.pool.QueryRow(ctx, "SELECT id, name, email, created_at, role from userprofile WHERE email=$1", userEmail)

	if err := row.Scan(&userProfile.ID, &userProfile.Name, &userProfile.Email, &userProfile.CreatedAt, &userProfile.Role); err != nil {
		return userProfile, err
	}

	return userProfile, nil
}
