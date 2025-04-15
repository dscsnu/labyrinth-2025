package controllers

import (
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

	rtr.HandleFunc("/api/eventlistener", Get(TeamChannelEventHandler(rtr)))

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

		//userEmail := r.Context().Value("email").(string)

		query := r.URL.Query()
		teamId := query.Get("team_id")

		//user, err := rtr.State.DB.GetUser(context.Background(), userEmail)
		//if err != nil {
		//
		//	http.Error(w, "internal error occurred", http.StatusInternalServerError)
		//	rtr.Logger.Error("internal error", "error", err)
		//	return
		//}

		//team, err := rtr.State.DB.GetTeamByUserId(context.Background(), user.ID)

		//if err != nil {

		//	http.Error(w, "internal error occurred, is the user in a valid team?", http.StatusInternalServerError)
		//	return
		//}

		teamChannel := rtr.State.ChanPool.GetChannel(teamId)
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

		//for eventMessage := range listenerChannel {

		//	if err := json.NewEncoder(w).Encode(eventMessage); err != nil {

		//		rtr.Logger.Debug("http stream write failed")

		//	}
		//	flusher.Flush()
		//}

		for eventMessage := range listenerChannel {
			// Convert event message to JSON string
			data, err := json.Marshal(eventMessage)
			if err != nil {
				rtr.Logger.Debug("failed to marshal event message", slog.String("error", err.Error()))
				continue
			}

			// Construct SSE message
			sseMessage := "data: " + string(data) + "\n\n"

			// Write the SSE message to the response
			_, err = w.Write([]byte(sseMessage))
			if err != nil {
				rtr.Logger.Debug("http stream write failed", slog.String("error", err.Error()))
				return
			}

			// Flush the output to ensure immediate delivery of the message
			flusher.Flush()
		}

	})

}
