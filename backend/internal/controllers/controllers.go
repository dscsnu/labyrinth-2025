package controllers

import (
	"context"
	"encoding/json"
	"labyrinth/internal/controllers/middleware"
	"labyrinth/internal/protocol"
	"labyrinth/internal/router"
	"log/slog"
	"net/http"

	_ "labyrinth/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func HandleAll(rtr *router.Router) {
	// GET Routes here
	rtr.HandleFunc("/api", Get(DefaultHandler(rtr)))
	rtr.HandleFunc("/api/team", middleware.Authorized(rtr, Get(GetTeamHandler(rtr))))
	rtr.HandleFunc("/api/game", Get(GameConfigHandler(rtr)))

	// POST Routes
	rtr.HandleFunc("/api/user/status", middleware.Authorized(rtr, Post(TeamMemberStatusUpdateHandler(rtr))))
	rtr.HandleFunc("/api/team/create", middleware.Authorized(rtr, Post(TeamCreationHandler(rtr))))
	rtr.HandleFunc("/api/team/update", middleware.Authorized(rtr, Post(TeamUpdateHandler(rtr))))
	rtr.HandleFunc("/api/team/leave", middleware.Authorized(rtr, Post(LeaveTeamHandler(rtr))))

	rtr.HandleFunc("/api/eventlistener", middleware.Authorized(rtr, TeamChannelEventHandler(rtr)))

	rtr.Handle("/swagger/", http.StripPrefix("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3100/swagger/doc.json"),
	)))
}

func DefaultHandler(rtr *router.Router) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if _, err := w.Write([]byte("/api is up")); err != nil {

			rtr.Logger.Error("error serving /api", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

		}

	})

}

func TeamChannelEventHandler(rtr *router.Router) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userEmail := r.Context().Value("email").(string)

		user, err := rtr.State.DB.GetUser(context.Background(), userEmail)
		if err != nil {

			http.Error(w, "internal error occurred", http.StatusInternalServerError)
			return
		}

		team, err := rtr.State.DB.GetTeamByUserId(context.Background(), user.ID)

		if err != nil {

			http.Error(w, "internal error occurred, is the user in a valid team?", http.StatusInternalServerError)
			return
		}

		teamChannel := rtr.State.ChanPool.GetChannel(team.ID)
		listenerChannel := make(chan protocol.Packet)
		teamChannel.AddMember(listenerChannel)

		//w.Header().Add("Content-Type", "text/event-stream")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		flusher, ok := w.(http.Flusher)
		if !ok {

			http.Error(w, "Could not create flusher", http.StatusInternalServerError)
			return

		}

		for eventMessage := range listenerChannel {

			if err = json.NewEncoder(w).Encode(eventMessage); err != nil {

				rtr.Logger.Debug("http stream write failed")

			}
			flusher.Flush()
		}

	})

}
