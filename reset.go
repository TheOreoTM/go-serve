package main

import "net/http"

func (cfg *WebServerContext) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.ApiConfig.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
