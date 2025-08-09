package cache

import (
	"context"
	"labyrinth/internal/database"
	"labyrinth/internal/types"
	"time"

	"github.com/google/uuid"
)

func (CM *CacheManager) GetUserByEmailCache(ctx context.Context, db *database.PostgresDriver, email string) (types.UserProfile, error) {
	var user types.UserProfile

	if cached, exists := CM.Get(UserProfile, email); exists {
		if typeCastUser, ok := cached.(types.UserProfile); ok {
			return typeCastUser, nil
		}
	}

	user, err := db.GetUser(ctx, email)
	if err != nil {
		return user, err
	}

	CM.Set(UserProfile, email, user, 60*time.Minute)
	CM.Set(UserProfile, user.ID.String(), user, 60*time.Minute)

	return user, nil
}

func (CM *CacheManager) GetUserByIdCache(ctx context.Context, db *database.PostgresDriver, userID string) (types.UserProfile, error) {
	var user types.UserProfile

	if cached, exists := CM.Get(UserProfile, userID); exists {
		if typeCastUser, ok := cached.(types.UserProfile); ok {
			return typeCastUser, nil
		}
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return user, err
	}

	user, err = db.GetUserById(ctx, userUUID)
	if err != nil {
		return user, err
	}

	CM.Set(UserProfile, userID, user, 60*time.Minute)
	CM.Set(UserProfile, user.Email, user, 60*time.Minute)

	return user, nil
}

func (CM *CacheManager) GetUserByUUIDCache(ctx context.Context, db *database.PostgresDriver, userID uuid.UUID) (types.UserProfile, error) {
	return CM.GetUserByIdCache(ctx, db, userID.String())
}

func (CM *CacheManager) InvalidateUserCacheByEmail(email string) {
	CM.Delete(UserProfile, email)
}

func (CM *CacheManager) InvalidateUserCacheByID(userID string) {
	CM.Delete(UserProfile, userID)
}

func (CM *CacheManager) UpdateUserReadyStateCache(ctx context.Context, db *database.PostgresDriver, userEmail string, status bool) error {
	err := db.UpdateUserReadyState(ctx, userEmail, status)
	if err != nil {
		return err
	}

	CM.InvalidateUserCacheByEmail(userEmail)

	return nil
}
