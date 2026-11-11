package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"weaver-edge/internal/transport/dto"
)

func MiddlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(secret string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			unauthorized := func(message string) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
					Error: message,
				})
			}

			authHeader := r.Header.Get("Authorization")
			parts := strings.Fields(authHeader)

			if len(parts) != 2 {
				unauthorized("invalid authorization header")
				return
			}

			if parts[0] != "Bearer" {
				unauthorized("invalid auth scheme")
				return
			}

			if parts[1] != secret {
				unauthorized("invalid token")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
