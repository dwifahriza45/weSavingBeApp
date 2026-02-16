package middleware

import (
	"BE_WE_SAVING/internal/shared/jwt"
	"BE_WE_SAVING/internal/shared/utils"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.JSONError(w, http.StatusUnauthorized, "NOK", "Missing Authorization Header", true)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.JSONError(w, http.StatusUnauthorized, "NOK", "Invalid Authorization Format", true)
				return
			}

			userID, err := jwt.ParseJWT(secret, parts[1])
			if err != nil {
				utils.JSONError(w, http.StatusUnauthorized, "NOK", "Invalid or Expired Token", true)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	return userID, ok
}
