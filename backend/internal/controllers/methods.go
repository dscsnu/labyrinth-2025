package controllers

import "net/http"

func Get(next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			return
		}

		next.ServeHTTP(w, r)

	}

}

func Post(next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			return
		}

		next.ServeHTTP(w, r)

	}

}

func Patch(next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPatch {
			return
		}

		next.ServeHTTP(w, r)

	}

}
