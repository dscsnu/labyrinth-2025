package controllers

import (
	"labyrinth/internal/router"
	"net/http"
)

func TeamCreationHandler(rtr *router.Router) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userEmail := r.Context().Value("email").(string)

	})

}
