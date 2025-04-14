package controllers

import (
	"labyrinth/internal/controllers/middleware"
	"labyrinth/internal/router"
	"log/slog"
	"net/http"

	_ "labyrinth/docs"

	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func HandleAll(rtr *router.Router) {

	cors.AllowAll()

	// GET Routes here
	rtr.HandleFunc("/api", Get(DefaultHandler(rtr)))
	rtr.HandleFunc("/api/team", middleware.Authorized(rtr, Get(GetTeamHandler(rtr))))
	rtr.HandleFunc("/api/game", Get(GameConfigHandler(rtr)))

	// POST Routes
	rtr.HandleFunc("/api/user/status", middleware.Authorized(rtr, Post(TeamMemberStatusUpdateHandler(rtr))))
	rtr.HandleFunc("/api/createteam", middleware.Authorized(rtr, Post(TeamCreationHandler(rtr))))
	rtr.HandleFunc("/api/updateteam", middleware.Authorized(rtr, Post(TeamUpdateHandler(rtr))))

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

		// get team id through db calls
		// teamChannel := rtr.State.ChanPool.GetChannel(teamId)
		// listenerChannel := make(chan protocol.Packet)
		// teamChannel.AddMember(listenerChannel)

		// flusher, ok := w.(http.Flusher)
		// if !ok {

		// 	http.Error(w, "Could not create flusher", http.StatusInternalServerError)
		// 	return

		// }

		// for eventMessage := range listenerChannel {

		// 	w.Write(something)
		// 	flusher.Flush()
		// }

	})

}
