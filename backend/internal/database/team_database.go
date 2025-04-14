package database

import (
	"context"
	"errors"
	"labyrinth/internal/types"
	"math/rand"
	"time"

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

	rows, err := pd.pool.Query(ctx, `
		SELECT u.id, u.name, u.email, u.created_at, u.role, tm.is_ready
		FROM teammember tm
		JOIN userprofile u ON u.id = tm.user_id
		WHERE tm.team_id = $1
		`, teamId)
	if err != nil {
		return team, err
	}
	defer rows.Close()

	var members []types.TeamMember

	for rows.Next() {
		var member types.TeamMember
		if err := rows.Scan(&member.ID, &member.Name, &member.Email, &member.CreatedAt, &member.Role, &member.IsReady); err != nil {
			return team, err
		}
		members = append(members, member)
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

func (pd *PostgresDriver) AssignLevelsToTeam(ctx context.Context, teamId string) error {

	levelIds := shuffleLevelIds()

	tx, err := pd.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	for seq, levelId := range levelIds {
		_, err := tx.Exec(ctx, `
			INSERT INTO teamlevelassignment (team_id, level_id, sequence)
			VALUES ($1, $2, $3)
		`, teamId, levelId, seq)

		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)

}

func shuffleLevelIds() []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ids := []int{1, 2, 3, 4, 5, 6}
	r.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})
	return append(ids, 7)
}
