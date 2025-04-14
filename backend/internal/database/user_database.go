package database

import (
	"context"
	"labyrinth/internal/types"

	"github.com/google/uuid"
)

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

func (pd *PostgresDriver) UpdateUserReadyState(ctx context.Context, userEmail string, status bool) error {

	user, err := pd.GetUser(ctx, userEmail)
	if err != nil {
		return err
	}

	_, err = pd.pool.Exec(ctx, "UPDATE teammember SET is_ready=$1 WHERE user_id=$2", status, user.ID)
	if err != nil {
		return err
	}

	return nil
}
