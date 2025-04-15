package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"labyrinth/internal/channel"
	"labyrinth/internal/protocol"
	"labyrinth/internal/router"
	"labyrinth/internal/types"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// TeamCreationHandler creates a new team and assigns default levels.
//
//	@Summary		Create Team
//	@Description	Creates a new team using the provided team name and returns the generated team ID. Also assigns default levels and initializes a communication channel.
//	@Tags			Team
//	@Accept			json
//	@Produce		json
//	@Param			body	body		object{team_name=string}	true	"Payload containing the team name"
//	@Success		200		{object}	object{team_id=string}		"The generated team ID"
//	@Failure		400		{object}	object{error=string}		"Bad Request"
//	@Failure		500		{object}	object{error=string}		"Internal Server Error"
//	@Security		BearerAuth
//	@Router			/api/createteam [post]
func TeamCreationHandler(rtr *router.Router) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userEmail := r.Context().Value("email").(string)

		t := struct {
			TeamName string `json:"team_name"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {

			//http.Error(w, "error reading teamName field, invalid json payload", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			apiResponse := types.ApiResponse{
				Success: false,
				Message: "invalid json payload",
				Payload: nil,
			}
			//json.NewEncoder(w).Encode(map[string]string{
			//	"error": "error reading tean_name field, invalid json payload",
			//})
			json.NewEncoder(w).Encode(apiResponse)
			return

		}

		profile, err := rtr.State.DB.GetUser(context.Background(), userEmail)
		if err != nil {

			//http.Error(w, "internal error occurred", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "internal server error",
			})
			rtr.Logger.Error("internal error", "error ", err.Error())
			return

		}

		teamId, err := rtr.State.DB.CreateTeam(context.Background(), t.TeamName, profile.ID)
		if err != nil {
			//http.Error(w, "error creating team in database", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "internal server error",
			})
			rtr.Logger.Error("internal error creating team", "error", err.Error())
			return
		}

		err = rtr.State.DB.AssignLevelsToTeam(context.Background(), teamId)
		if err != nil {
			//http.Error(w, "error assigning levels to team", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "internal server error",
			})
			rtr.Logger.Error("internal error assigning levels", "error", err.Error())
			return
		}

		teamChannel := channel.NewChannel()
		rtr.State.ChanPool.AddChannel(teamId, teamChannel)

		go teamChannel.Start()

		json.NewEncoder(w).Encode(map[string]string{
			"team_id": teamId,
		})

	})

}

// TeamUpdateHandler adds members to a specified team
//
//	@Summary		Add member to team
//	@Description	Adds members to the team specified in the payload
//	@Tags			Team
//	@Accept			json
//	@Produce		json
//	@Param			body	body		object{team_id=string}	true	"Payload containing the id of the team to join"
//	@Success		200		{object}	types.Team				"Updated team the member is added to"
//	@Failure		400		{object}	object{error=string}	"Bad request"
//	@Failure		500		{object}	object{error=string}	"Internal Server Error"
//	@Router			/api/updateteam [post]
func TeamUpdateHandler(rtr *router.Router) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userEmail := r.Context().Value("email").(string)

		t := struct {
			TeamId string `json:"team_id"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {

			http.Error(w, "error getting team_id field, invalid json payload", http.StatusBadRequest)
			return

		}

		profile, err := rtr.State.DB.GetUser(context.Background(), userEmail)

		if err != nil {

			http.Error(w, "internal error occurred", http.StatusInternalServerError)
			rtr.Logger.Error("internal error", "error", err.Error())
			return

		}

		if err := rtr.State.DB.AddTeamMember(context.Background(), t.TeamId, profile.ID); err != nil {
			http.Error(w, "the team is full", http.StatusInternalServerError)
			return
		}

		team, err := rtr.State.DB.GetTeamByID(context.Background(), t.TeamId)
		if err != nil {
			http.Error(w, "error fetching the team", http.StatusInternalServerError)
			rtr.Logger.Error("internal error while getting team", "error", err.Error())
		}

		teamChannel := rtr.State.ChanPool.GetChannel(team.ID)
		teamChannel.Broadcast(protocol.Packet{Type: "BackgroundMessage", BackgroundMessage: protocol.BackgroundMessage{Relay: fmt.Sprintf("teamId:%s -> %s joined the team", team.ID, profile.Email), MsgContext: "channel_creation"}})

		if err := json.NewEncoder(w).Encode(team); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
		}

	})

}

// GetTeamHandler returns the team from a given team ID or for a given user ID
//
//	@Summary		Get team info
//	@Description	Gets the team info, using either team ID or user ID
//	@Tags			Team
//	@Accept			json
//	@Produce		json
//	@Param			team_id	query		string					false	"ID of the team"
//	@Param			user_id	query		string					false	"ID of a user belonging to the team"
//	@Success		200		{object}	types.Team				"Team retrieved successfully"
//	@Failure		400		{object}	object{error=string}	"Bad request"
//	@Failure		500		{object}	object{error=string}	"Internal server error"
//	@Router			/api/tean [get]
func GetTeamHandler(rtr *router.Router) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		teamId := query.Get("team_id")
		userId := query.Get("user_id")

		if teamId == "" && userId == "" {
			//http.Error(w, "Either user_id or team_id must be provided", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "either used_id or team_id not provided",
			})
			return
		}

		var team = types.Team{}
		var err error

		if teamId != "" {
			team, err = rtr.State.DB.GetTeamByID(context.Background(), teamId)
		} else if userId != "" {
			parsedId, parseErr := uuid.Parse(userId)
			if parseErr != nil {
				//http.Error(w, "invalid player_id format", http.StatusBadRequest)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "player id format is invalid, should be uuid string",
				})
			}
			team, err = rtr.State.DB.GetTeamByUserId(context.Background(), parsedId)
		}

		if err != nil {
			rtr.Logger.Error("failed to fetch team", slog.String("error", err.Error()))
			//http.Error(w, "internal server error", http.StatusInternalServerError)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "internal server error",
			})
			return
		}

		if err := json.NewEncoder(w).Encode(team); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "internal server error",
			})
			//http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})

}
