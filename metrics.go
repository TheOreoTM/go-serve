package main

import (
	"fmt"
	"net/http"
)

// func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Add("Content-Type", "text/html")
// 	w.WriteHeader(http.StatusOK)
// }

func (cfg *WebServerContext) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.ApiConfig.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *WebServerContext) handlerAdminMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `
	<html>
    <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited %d times!</p>
    </body>
  </html>`,
		cfg.ApiConfig.fileserverHits)
}
