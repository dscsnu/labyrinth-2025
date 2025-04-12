package controllers

import (
	"context"
	"encoding/json"
	"labyrinth/internal/router"
	"net/http"
)

func TeamMemberStatusUpdateHandler(rtr *router.Router) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userEmail := r.Context().Value("email").(string)

		payload := struct {
			UserStatus bool `json:"user_status"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "error reading userStatus field, invalid json payload", http.StatusBadRequest)
			return
		}

		err := rtr.State.DB.UpdateUserReadyState(context.Background(), userEmail, payload.UserStatus)
		if err != nil {
			http.Error(w, "internal error occurred", http.StatusInternalServerError)
			rtr.Logger.Error("internal error", "error", err.Error())
			return
		}

		user, err := rtr.State.DB.GetUser(context.Background(), userEmail)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			rtr.Logger.Error("internal error", "error", err.Error())
			return
		}

		team, err := rtr.State.DB.GetTeamByUserId(context.Background(), user.ID)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			rtr.Logger.Error("internal error", "error", err.Error())
			return
		}

		err = json.NewEncoder(w).Encode(team)
		if err != nil {
			http.Error(w, "error encoding json", http.StatusInternalServerError)
			rtr.Logger.Error("internal error", "error", err.Error())
			return
		}
	})

}
