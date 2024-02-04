package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AvivKermann/BlogAgg/internal/database"
	"github.com/AvivKermann/BlogAgg/internal/jsonResponse"
)

func (cfg *apiConfig) handlerGetPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	strLimit := r.URL.Query().Get("limit")
	fmt.Println(strLimit)
	limit, err := strconv.Atoi(strLimit)
	if err != nil {
		jsonResponse.RespondWithError(w, http.StatusBadRequest, "invalid limit")
		return
	}
	dbPosts, err := cfg.DB.GetPostsByUserID(r.Context(), database.GetPostsByUserIDParams{
		ID:    user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		jsonResponse.RespondWithError(w, http.StatusBadRequest, "user not found")
		return
	}
	posts := []Post{}
	for _, dbPost := range dbPosts {
		posts = append(posts, databasePostToPost(dbPost))
	}
	jsonResponse.ResponsWithJson(w, http.StatusOK, posts)
}
