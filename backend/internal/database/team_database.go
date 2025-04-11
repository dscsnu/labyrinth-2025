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
		return errors.New("fullteam")
	}

	if _, err := pd.pool.Exec(ctx, "INSERT INTO teammember(team_id, user_id) VALUES ($1,$2)", teamId, userId); err != nil {
		return err
	}

	return nil
}

func (pd *PostgresDriver) GetTeamByID(ctx context.Context, teamId string) (types.Team, error) {
	team := types.Team{}

	err := pd.pool.QueryRow(ctx, "SELECT id, name FROM team WHERE team.id=$1", teamId).Scan(&team.ID, &team.Name)
	if err != nil {
		return team, err
	}

	rows, err := pd.pool.Query(ctx, "SELECT user_id FROM teammember WHERE team_id=$1", teamId)
	if err != nil {
		return team, err
	}
	defer rows.Close()

	var members []types.UserProfile

	for rows.Next() {
		var userId uuid.UUID

		if err := rows.Scan(&userId); err != nil {
			return team, err
		}

		user, err := pd.GetUserById(ctx, userId)
		if err != nil {
			return team, err
		}

		members = append(members, user)
	}

	team.Members = members
	return team, nil

}

func (pd *PostgresDriver) GetTeamByUserId(ctx context.Context, userId uuid.UUID) (types.Team, error) {

	var team = types.Team{}

	err := pd.pool.QueryRow(ctx, "SELECT team_id FROM teammember WHERE teammember.user_id=$1", userId).Scan(&team.ID)
	if err != nil {
		return team, err
	}

	team, err = pd.GetTeamByID(ctx, team.ID)
	if err != nil {
		return team, err
	}

	return team, nil

}

func (pd *PostgresDriver) GetUser(ctx context.Context, userEmail string) (types.UserProfile, error) {
	userProfile := types.UserProfile{}
	row := pd.pool.QueryRow(ctx, "SELECT id, name, email, created_at, role from userprofile WHERE email=$1", userEmail)

	if err := row.Scan(&userProfile.ID, &userProfile.Name, &userProfile.Email, &userProfile.CreatedAt, &userProfile.Role); err != nil {
		return userProfile, err
	}

	return userProfile, nil
}

func (pd *PostgresDriver) GetUserById(ctx context.Context, userId uuid.UUID) (types.UserProfile, error) {
	user := types.UserProfile{}
	row := pd.pool.QueryRow(ctx, "SELECT id, name, email, created_at, role FROM userprofile WHERE id=$1", userId)

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.Role); err != nil {
		return user, err
	}

	return user, nil
}
