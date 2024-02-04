package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AvivKermann/BlogAgg/internal/database"
	"github.com/AvivKermann/BlogAgg/internal/jsonResponse"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	type paramaters struct {
		FeedId string `json:"feed_id"`
	}

	params := paramaters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		jsonResponse.RespondWithError(w, http.StatusBadRequest, "invalid feed id")
		return
	}
	feedId, err := uuid.Parse(params.FeedId)
	if err != nil {
		jsonResponse.RespondWithError(w, http.StatusBadRequest, "invalid feed id")
		return
	}

	dbfeedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feedId,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	feedFollow := databaseFeedFollowToFeedFollow(dbfeedFollow)
	if err != nil {
		jsonResponse.RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	jsonResponse.ResponsWithJson(w, http.StatusOK, feedFollow)

}
