package controllers

import (
	"labyrinth/internal/controllers/middleware"
	"labyrinth/internal/router"
	"log/slog"
	"net/http"
)

func HandleAll(rtr *router.Router) {
	// GET Routes here
	rtr.HandleFunc("/api", Get(DefaultHandler(rtr)))
	rtr.HandleFunc("/api/createteam", middleware.Authorized(rtr, Post(TeamCreationHandler(rtr))))
	rtr.HandleFunc("/api/updateteam", middleware.Authorized(rtr, Post(TeamUpdateHandler(rtr))))
}

func DefaultHandler(rtr *router.Router) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if _, err := w.Write([]byte("/api is up")); err != nil {

			rtr.Logger.Error("error serving /api", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

		}

	})

}
