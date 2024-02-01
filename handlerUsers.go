package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/AvivKermann/BlogAgg/internal/database"
	"github.com/AvivKermann/BlogAgg/internal/jsonResponse"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type paramaters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := paramaters{}
	err := decoder.Decode(&params)
	if err != nil {
		jsonResponse.RespondWithError(w, http.StatusBadRequest, "invalid user information")
		return
	}
	uuid, err := uuid.NewRandom()
	if err != nil {
		jsonResponse.RespondWithError(w, http.StatusInternalServerError, "cannot create new uuid")
		return
	}
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		jsonResponse.RespondWithError(w, http.StatusBadRequest, "cannot create user")
		return
	}
	jsonResponse.ResponsWithJson(w, http.StatusOK, user)

}
func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	jsonResponse.ResponsWithJson(w, http.StatusOK, databaseUserToUser(user))

}

func stripPrefix(key string) string {
	token := strings.TrimPrefix(key, "ApiKey ")
	return token
}
