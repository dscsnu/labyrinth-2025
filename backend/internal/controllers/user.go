package controllers

import (
	"context"
	"encoding/json"
	"labyrinth/internal/cache"
	queue "labyrinth/internal/controllers/db_queue"
	"labyrinth/internal/router"
	"labyrinth/internal/types"
	"net/http"
	"time"
)

// TeamMemberStatusUpdateHandler handles modifying the is_ready state for users in a team
//
//	@Summary		Modify is_ready state for members
//	@Description	Changes is_ready status for a user if they're in a team
//	@Tags			Team
//	@Accept			json
//	@Produce		json
//	@Param			body	body		object{user_status=bool}	true	"The ready state to change to"
//	@Success		200		{object}	types.Team				"The team with the updated ready state for the member"
//	@Failure		400		{object}	object{error=string}	"Bad request"
//	@Failure		500		{object}	object{error=string}	"Internal server error"
//	@Router			/api/user/status [post]
func TeamMemberStatusUpdateHandler(rtr *router.Router) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userEmail := r.Context().Value("email").(string)

		payload := struct {
			UserStatus bool `json:"user_status"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			apiResponse := types.ApiResponse{
				Success: false,
				Message: "invalid json payload",
				Payload: nil,
			}
			json.NewEncoder(w).Encode(apiResponse)
			return
		}

		// Get user profile from cache (fallback to DB)
		user, err := rtr.State.CM.GetUserByEmailCache(context.Background(), rtr.State.DB, userEmail)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to get user profile",
			})
			rtr.Logger.Error("failed to get user from cache", "error", err.Error())
			return
		}

		// Get team from cache (fallback to DB)
		team, err := rtr.State.CM.GetTeamByUserIdCache(context.Background(), rtr.State.DB, user.ID.String())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to get team",
			})
			rtr.Logger.Error("failed to get team from cache", "error", err.Error())
			return
		}

		// Update ready status in cache instantly
		updatedTeam := team
		for i, member := range updatedTeam.Members {
			if member.ID == user.ID {
				updatedTeam.Members[i].IsReady = payload.UserStatus
				break
			}
		}
		rtr.State.CM.Set(cache.Team, updatedTeam.ID, updatedTeam, 30*time.Minute)
		rtr.State.CM.Set(cache.TeamUserIDIndex, user.ID.String(), updatedTeam, 30*time.Minute)

		// Respond immediately with updated cache
		responsePayload, _ := json.Marshal(updatedTeam)
		apiResponse := types.ApiResponse{
			Success: true,
			Message: "Successfully updated ready status!",
			Payload: responsePayload,
		}
		json.NewEncoder(w).Encode(apiResponse)

		// Queue DB update
		queueReq := queue.DBQueueRequest{
			Type: "user_ready_update",
			Payload: map[string]interface{}{
				"user_email":  userEmail,
				"user_id":     user.ID,
				"team_id":     updatedTeam.ID,
				"user_status": payload.UserStatus,
			},
			Handler: func() error {
				err := rtr.State.DB.UpdateUserReadyState(context.Background(), userEmail, payload.UserStatus)
				if err != nil {
					rtr.Logger.Error("database error updating ready state", "error", err.Error())
					return err
				}
				// Optionally, refresh cache from DB after DB update
				user, err := rtr.State.DB.GetUser(context.Background(), userEmail)
				if err == nil {
					rtr.State.CM.Set(cache.UserProfile, userEmail, user, 60*time.Minute)
					rtr.State.CM.Set(cache.UserProfile, user.ID.String(), user, 60*time.Minute)
				}
				return nil
			},
		}
		queue.AddToQueue(queueReq)
	})

}
