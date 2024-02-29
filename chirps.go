package main

import "net/http"

func (ctx *WebServerContext) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := ctx.Database.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps")
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
