package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (ctx *WebServerContext) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := ctx.Database.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps")
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (ctx *WebServerContext) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := ctx.Database.GetChirp(idInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirp")
		return
	}

	if chirp.ID == 0 {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (ctx *WebServerContext) handlePostChirp(w http.ResponseWriter, r *http.Request) {
	var parameters struct {
		Body string `json:"body"`
	}

	w.Header().Add("Content-Type", "application/json")

	err := DecodeJSONBody(r, &parameters)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if parameters.Body == "" {
		respondWithError(w, 400, "Chirp cannot be empty")
		return
	}

	if len(parameters.Body) > 140 {
		respondWithError(w, 400, "Chirp cannot be longer than 140 characters")
		return
	}

	chirp, err := ctx.Database.CreateChirp(parameters.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error adding chirp")
		return
	}

	response := fmt.Sprintf(`{"id": "%d", "body": "%s"}`, chirp.ID, chirp.Body)

	respondWithJSON(w, http.StatusCreated, response)
}

func censorProfane(str string) string {
	var cleaned []string
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(str, " ")
	for _, word := range words {
		for _, profane := range profaneWords {
			if strings.ToLower(word) == profane {
				word = "****"
			}
		}
		cleaned = append(cleaned, word)
	}

	return strings.Join(cleaned, " ")
}
