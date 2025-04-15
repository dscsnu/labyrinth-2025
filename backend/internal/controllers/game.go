package controllers

import (
	"context"
	"encoding/json"
	"labyrinth/internal/router"
	"labyrinth/internal/types"
	"net/http"
)

// GameConfigHandler gets information about the game
//
//	@Summary		Get GameConfig
//	@Description	Gets the config data for the game
//	@Tags			Game
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string		"Game config retrieved successfully"
//	@Failure		500	{object}	object{error=string}	"Internal server error"
//	@Router			/api/game [get]
func GameConfigHandler(rtr *router.Router) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		gameConfig, err := rtr.State.DB.GetGameConfig(context.Background())
		if err != nil {
			//http.Error(w, "internal server error", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "internal server error",
			})
			rtr.Logger.Error("internal error getting gameconfig", "error", err)
			return
		}

		responsePayload, _ := json.Marshal(gameConfig)
		apiResponse := types.ApiResponse{
			Success: true,
			Message: "",
			Payload: responsePayload,
		}

		err = json.NewEncoder(w).Encode(apiResponse)
		if err != nil {
			//http.Error(w, "internal server error", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "internal server error",
			})
			rtr.Logger.Error("error encoding json response", "error", err)
			return
		}

	})

}
