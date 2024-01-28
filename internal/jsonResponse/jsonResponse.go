package jsonResponse

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(ErrorResponse{
		Error: message,
	})
	if err != nil {
		log.Printf("error marshalling JSON %s", err)
		statusCode = http.StatusInternalServerError
		return
	}
	w.WriteHeader(statusCode)
	w.Write(data)
}

func ResponsWithJson(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(statusCode)
	w.Write(data)
}
