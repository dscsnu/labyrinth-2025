package middleware

import (
	"context"
	"fmt"
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
			
			rtr.Logger.Error("JWT Token not found")

		}
		
		jwtToken, err := jwtauth.VerifyToken(attemptToken, rtr.State.JwtSession.SecretKey)
		if err != nil {
			fmt.Println(err)
			
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
func CORSFunc(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "*") // For development
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // Handle preflight requests
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        // Call the next handler
        next(w, r)
    }
}