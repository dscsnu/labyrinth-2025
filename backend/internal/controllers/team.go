package controllers

import (
	"context"
	"encoding/json"
	"labyrinth/internal/router"
	"net/http"
)

func TeamCreationHandler(rtr *router.Router) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userEmail := r.Context().Value("email").(string)

		t := struct {
			TeamName string `json:"teamName"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {

			http.Error(w, "error reading teamName field, invalid json payload", http.StatusBadRequest)
			return

		}

		profile, err := rtr.State.DB.GetUser(context.Background(), userEmail)
		if err != nil {

			http.Error(w, "internal error occurred", http.StatusInternalServerError)
			return

		}

		rtr.State.DB.CreateTeam(context.Background(), t.TeamName, profile.Email)

	})

}
