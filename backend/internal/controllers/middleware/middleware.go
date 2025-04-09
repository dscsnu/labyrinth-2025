package middleware

import (
	"context"
	"labyrinth/internal/jwtauth"
	"labyrinth/internal/router"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func getToken(r *http.Request) (string, bool) {

	authHeader := strings.Split(r.Header.Get("Authorization"), " ")
	if authHeader[0] != "Bearer" {

		return "", false

	}

	return authHeader[1], true

}

func Authorized(rtr *router.Router, next http.Handler) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		attemptToken, ok := getToken(r)
		if !ok {

			return

		}

		jwtToken, err := jwtauth.VerifyToken(attemptToken, rtr.State.JwtSession.SecretKey)
		if err != nil {

			return

		}

		mapclaims := jwtToken.Claims.(jwt.MapClaims)
		email, ok := mapclaims["email"]

		if !ok {
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(context.Background(), "email", email)))

	})

}
