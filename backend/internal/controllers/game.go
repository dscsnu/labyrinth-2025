package controllers

import (
	"context"
	"encoding/json"
	"labyrinth/internal/router"
	"net/http"
)

func GameConfigHandler(rtr *router.Router) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		gameConfig, err := rtr.State.DB.GetGameConfig(context.Background())
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			rtr.Logger.Error("internal error getting gameconfig", "error", err)
			return
		}

		err = json.NewEncoder(w).Encode(gameConfig)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			rtr.Logger.Error("error encoding json response", "error", err)
			return
		}

	})

}
