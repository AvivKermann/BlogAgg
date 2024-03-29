package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AvivKermann/BlogAgg/internal/database"
	"github.com/AvivKermann/BlogAgg/internal/jsonResponse"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type paramaters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := paramaters{}
	err := decoder.Decode(&params)

	if err != nil {
		jsonResponse.RespondWithError(w, http.StatusBadRequest, "cannot decode feed parameters")
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		jsonResponse.RespondWithError(w, http.StatusBadRequest, "cannot create feed")
		return
	}
	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feed.ID,
		UserID:    feed.UserID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		jsonResponse.RespondWithError(w, http.StatusBadRequest, "cannot create feed follow")
		return
	}
	jsonResponse.ResponsWithJson(w, http.StatusOK, struct {
		Feed        Feed       `json:"feed"`
		Feed_follow FeedFollow `json:"feed_follow"`
	}{
		Feed:        databaseFeedToFeed(feed),
		Feed_follow: databaseFeedFollowToFeedFollow(feedFollow),
	})
}

func (cfg *apiConfig) handlerGetAllFeeds(w http.ResponseWriter, r *http.Request) {

	dbFeeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		jsonResponse.RespondWithError(w, http.StatusBadRequest, "cannot get feeds")
		return
	}
	feeds := []Feed{}
	for _, dbFeed := range dbFeeds {
		feeds = append(feeds, databaseFeedToFeed(dbFeed))
	}
	jsonResponse.ResponsWithJson(w, http.StatusOK, feeds)

}
