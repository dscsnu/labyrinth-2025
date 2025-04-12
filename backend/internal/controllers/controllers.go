package controllers

import (
	"labyrinth/internal/controllers/middleware"
	"labyrinth/internal/router"
	"log/slog"
	"net/http"

	"github.com/rs/cors"
)

func HandleAll(rtr *router.Router) {

	cors.AllowAll()

	// GET Routes here
	rtr.HandleFunc("/api", Get(DefaultHandler(rtr)))
	rtr.HandleFunc("/api/createteam", middleware.Authorized(rtr, Post(TeamCreationHandler(rtr))))
	rtr.HandleFunc("/api/updateteam", middleware.Authorized(rtr, Post(TeamUpdateHandler(rtr))))
	rtr.HandleFunc("/api/team", middleware.Authorized(rtr, Get(GetTeamHandler(rtr))))
	rtr.HandleFunc("/api/user/status", middleware.Authorized(rtr, Post(TeamMemberStatusUpdateHandler(rtr))))
}

func DefaultHandler(rtr *router.Router) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if _, err := w.Write([]byte("/api is up")); err != nil {

			rtr.Logger.Error("error serving /api", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

		}

	})

}
