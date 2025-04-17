package cache

import (
	"context"
	"labyrinth/internal/database"
	"labyrinth/internal/types"
	"time"

	"github.com/google/uuid"
)

func (CM *CacheManager) GetTeamByUserIdCache(ctx context.Context, db *database.PostgresDriver, userId string) (types.Team, error) {
	var team = types.Team{}
	var err = error(nil)

	parsedId, parseErr := uuid.Parse(userId)
	if parseErr != nil {
		//http.Error(w, "invalid player_id format", http.StatusBadRequest)
		//w.WriteHeader(http.StatusBadRequest)
		//json.NewEncoder(w).Encode(map[string]string{
		//	"error": "player id format is invalid, should be uuid string",
		//})

		//apiResponse := types.ApiResponse{
		//	Success: false,
		//	Message: "player id format is invalid, should be uuid string",
		//	Payload: nil,
		//}
		//json.NewEncoder(w).Encode(apiResponse)
		return team, err

	}

	if cached, exists := CM.Get(TeamUserIDIndex, userId); exists {
		if typeCastTeam, ok := cached.(types.Team); ok {
			return typeCastTeam, nil
		}
	} else {
		team, err = db.GetTeamByUserId(ctx, parsedId)
		if err != nil {
			return team, err
		}
		CM.Set(TeamUserIDIndex, userId, team, 30*time.Minute)
	}

	return team, err
}

func (CM *CacheManager) GetTeamByIdCache(ctx context.Context, db *database.PostgresDriver, teamId string) (types.Team, error) {

	var team = types.Team{}
	if cached, exists := CM.Get(Team, teamId); exists {
		if typeCastTeam, ok := cached.(types.Team); ok {
			return typeCastTeam, nil
		}
	} else {
		team, err := db.GetTeamByID(ctx, teamId)
		if err != nil {
			return team, err
		}
		CM.Set(Team, teamId, team, 60*time.Minute)
	}

	return team, nil

}
