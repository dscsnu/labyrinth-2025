package controllers

import (
	"context"
	"encoding/json"
	"labyrinth/internal/router"
	"net/http"
)

// TeamMemberStatusUpdateHandler handles modifying the is_ready state for users in a team
//
//	@Summary		Modify is_ready state for members
//	@Description	Changes is_ready status for a user if they're in a team
//	@Tags			Team
//	@Accept			json
//	@Produce		json
//	@Param			body	body		object{is_ready=bool}	true	"The ready state to change to"
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
			//http.Error(w, "error reading userStatus field, invalid json payload", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid json payload",
			})
			return
		}

		err := rtr.State.DB.UpdateUserReadyState(context.Background(), userEmail, payload.UserStatus)
		if err != nil {
			//http.Error(w, "internal error occurred", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "internal server error",
			})
			rtr.Logger.Error("internal error", "error", err.Error())
			return
		}

		user, err := rtr.State.DB.GetUser(context.Background(), userEmail)
		if err != nil {
			//http.Error(w, "internal server error", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "internal server error",
			})
			rtr.Logger.Error("internal error", "error", err.Error())
			return
		}

		team, err := rtr.State.DB.GetTeamByUserId(context.Background(), user.ID)
		if err != nil {
			//http.Error(w, "internal server error", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "internal server error",
			})
			rtr.Logger.Error("internal error", "error", err.Error())
			return
		}

		err = json.NewEncoder(w).Encode(team)
		if err != nil {
			//http.Error(w, "error encoding json", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "internal server error",
			})
			rtr.Logger.Error("internal error", "error", err.Error())
			return
		}
	})

}
