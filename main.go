package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	app_router := chi.NewRouter()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	app_router.Handle("/app", fsHandler)
	app_router.Handle("/app/*", fsHandler)

	api_router := chi.NewRouter()

	api_router.Get("/healthz", handlerReadiness)
	api_router.Get("/metrics", apiCfg.handlerMetrics)
	api_router.Get("/reset", apiCfg.handlerReset)

	app_router.Mount("/api", api_router)

	corsMux := middlewareCors(app_router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
