package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/theoreotm/go-serve/database"
)

type apiConfig struct {
	fileserverHits int
}

type WebServerContext struct {
	ApiConfig apiConfig
	Database  *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	database, err := database.NewDB("./database.json")
	if err != nil {
		panic(err)
	}

	context := WebServerContext{
		ApiConfig: apiConfig{
			fileserverHits: 0,
		},
		Database: database,
	}

	router := chi.NewRouter()
	fsHandler := context.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot+"/app"))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", context.handlerReset)
	apiRouter.Get("/reset", context.handleGetChirps)
	apiRouter.Post("/validate_chirp", handleValidateChirp)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", context.handlerAdminMetrics)

	router.Mount("/admin", adminRouter)
	router.Mount("/api", apiRouter)

	corsMux := middlewareCors(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
