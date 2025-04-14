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

func TeamCreationHandler(rtr *router.Router) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userEmail := r.Context().Value("email").(string)

		t := struct {
			TeamName string `json:"team_name"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {

			http.Error(w, "error reading teamName field, invalid json payload", http.StatusBadRequest)
			return

		}

		profile, err := rtr.State.DB.GetUser(context.Background(), userEmail)
		if err != nil {

			http.Error(w, "internal error occurred", http.StatusInternalServerError)
			rtr.Logger.Error("internal error", "error ", err.Error())
			return

		}

		teamId, err := rtr.State.DB.CreateTeam(context.Background(), t.TeamName, profile.ID)
		if err != nil {
			http.Error(w, "error creating team in database", http.StatusInternalServerError)
			rtr.Logger.Error("internal error creating team", "error", err.Error())
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

func GetTeamHandler(rtr *router.Router) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		teamId := query.Get("team_id")
		userId := query.Get("user_id")

		if teamId == "" && userId == "" {
			http.Error(w, "Either user_id or team_id must be provided", http.StatusBadRequest)
			return
		}

		var team = types.Team{}
		var err error

		if teamId != "" {
			team, err = rtr.State.DB.GetTeamByID(context.Background(), teamId)
		} else if userId != "" {
			parsedId, parseErr := uuid.Parse(userId)
			if parseErr != nil {
				http.Error(w, "invalid player_id format", http.StatusBadRequest)
			}
			team, err = rtr.State.DB.GetTeamByUserId(context.Background(), parsedId)
		}

		if err != nil {
			rtr.Logger.Error("failed to fetch team", slog.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(team); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})

}
