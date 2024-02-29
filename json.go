package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func DecodeJSONBody(r *http.Request, params interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(params)
	if err != nil {
		return fmt.Errorf("error decoding parameters: %w", err)
	}
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, msg)))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
