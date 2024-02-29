package main

import (
	"net/http"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	w.Header().Add("Content-Type", "application/json")

	var params parameters
	err := DecodeJSONBody(r, &params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if params.Body == "" {
		respondWithError(w, 400, "Chirp cannot be empty")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp cannot be longer than 140 characters")
		return
	}

	respondWithJSON(w, 200, `{"cleaned_body": "`+censorProfane(params.Body)+`"}`)
}
