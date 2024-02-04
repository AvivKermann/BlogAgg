package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AvivKermann/BlogAgg/internal/database"
	"github.com/AvivKermann/BlogAgg/internal/jsonResponse"
	"github.com/go-chi/chi/v5"
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

func (cfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	strFeedFollowID := chi.URLParam(r, "feedFollowID")
	feedFollowId, err := uuid.Parse(strFeedFollowID)
	if err != nil {
		jsonResponse.RespondWithError(w, http.StatusBadRequest, "invalid feed id")
		return
	}
	feedFollow, err := cfg.DB.GetFeedFollowByID(r.Context(), feedFollowId)
	if err != nil {
		jsonResponse.RespondWithError(w, http.StatusNotFound, "id dosent exist")
		return
	}
	if err := cfg.DB.DeleteFeedFollowByID(r.Context(), feedFollow.ID); err != nil {
		jsonResponse.RespondWithError(w, http.StatusNotFound, "id dosent exist")
		return
	}
	jsonResponse.ResponsWithJson(w, http.StatusOK, struct {
		Body string `json:"body"`
	}{
		Body: fmt.Sprintf("feed follow %v deleted", strFeedFollowID),
	})
}
func (cfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.GetAllFeedFollowsByUserID(r, user.ID)
	if err != nil {
		jsonResponse.RespondWithError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	jsonResponse.ResponsWithJson(w, http.StatusOK, feedFollows)
}

func (cfg *apiConfig) GetAllFeedFollowsByUserID(r *http.Request, userID uuid.UUID) ([]FeedFollow, error) {
	feedFollows := []FeedFollow{}
	dbFeedFollows, err := cfg.DB.GetFeedFollowByUserID(r.Context(), userID)
	if err != nil {
		return feedFollows, err
	}
	for _, dbFeedFollow := range dbFeedFollows {
		feedFollow := databaseFeedFollowToFeedFollow(dbFeedFollow)
		feedFollows = append(feedFollows, feedFollow)
	}
	return feedFollows, nil
}
