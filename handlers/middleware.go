package handlers

import (
	"log"
	"net/http"
	"strings"

	"go-backend-crm/config"
	"go-backend-crm/errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/cors"
)

// Middleware de autenticación JWT
func jwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "No autorizado", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(strings.TrimPrefix(tokenString, "Bearer "), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.NewCustomError("método de firma inesperado", http.StatusUnauthorized)
			}
			return []byte(config.GetConfig().JWTSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "No autorizado", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Middleware de CORS
func corsMiddleware(next http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	return c.Handler(next)
}

// Middleware de logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
