package main

import (
	"net/http"

	"github.com/AvivKermann/BlogAgg/internal/jsonResponse"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	jsonResponse.ResponsWithJson(w, http.StatusOK, struct {
		Body string `json:"status"`
	}{
		Body: "ok",
	})
}
func handlerErr(w http.ResponseWriter, r *http.Request) {
	jsonResponse.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
