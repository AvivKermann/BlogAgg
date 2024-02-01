package main

import (
	"context"
	"net/http"

	"github.com/AvivKermann/BlogAgg/internal/database"
	"github.com/go-chi/cors"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := stripPrefix(r.Header.Get("Authorization"))
		user, err := cfg.authenticateUser(apiKey, r.Context())
		if err != nil {
			http.Error(w, "Unauthorize", http.StatusUnauthorized)
		}
		handler(w, r, user)
	}

}

func (cfg *apiConfig) authenticateUser(apiKey string, context context.Context) (database.User, error) {
	user, err := cfg.DB.GetUserByApikey(context, apiKey)
	if err != nil {
		return database.User{}, err
	}
	return user, nil
}

var middleware = cors.Options{
	AllowedOrigins:   []string{"https://*", "http://*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: false,
}
