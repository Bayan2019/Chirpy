package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Bayan2019/chirpy/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load(".env")
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if dbg != nil && *dbg {
		err := db.ResetDB()
		if err != nil {
			log.Fatal(err)
		}
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
	}

	app_router := chi.NewRouter()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	app_router.Handle("/app", fsHandler)
	app_router.Handle("/app/*", fsHandler)

	api_router := chi.NewRouter()

	api_router.Get("/healthz", handlerReadiness)
	api_router.Get("/reset", apiCfg.handlerReset)

	api_router.Post("/revoke", apiCfg.handlerRevoke)
	api_router.Post("/refresh", apiCfg.handlerRefresh)
	api_router.Post("/login", apiCfg.handlerLogin)
	api_router.Post("/users", apiCfg.handlerUsersCreate)

	api_router.Put("/users", apiCfg.handlerUsersUpdate)

	api_router.Post("/chirps", apiCfg.handlerChirpsCreate)

	api_router.Get("/chirps", apiCfg.handlerChirpsRetrieve)
	api_router.Get("/chirps/{chirpID}", apiCfg.handlerChirpsGet)

	app_router.Mount("/api", api_router)

	admin_router := chi.NewRouter()

	admin_router.Get("/metrics", apiCfg.handlerMetrics)

	app_router.Mount("/admin", admin_router)

	corsMux := middlewareCors(app_router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
