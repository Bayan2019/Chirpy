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

// 2. Routing / 1. Stateful Handlers
// Create a struct that will hold any stateful, in-memory data we'll need to keep track of
//
//	type apiConfig struct {
//		fileserverHits int
//	}
type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
	polkaKey       string
}

func main() {
	// 1. Servers / 5. Fileservers
	// const filepathRoot = "."
	const filepathRoot = "."

	// 1. Servers / 4. Server
	// const port = "8080"
	const port = "8080"

	godotenv.Load(".env")
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY environment variable is not set")
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
		polkaKey:       polkaKey,
	}

	// 1. Servers / 4. Server
	// Create a new http.ServeMux
	// mux := http.NewServeMux()
	app_router := chi.NewRouter()

	// 1. Servers / 5. Fileservers
	// Use a standard http.FileServer as the handler
	// Use http.Dir to convert a filepath to a directory for the http.FileServer
	// mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	// 2. Routing / 1. Stateful Handlers
	// Wrap the http.FileServer handler with the middleware "middlewareMetricsInc"
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	// 1. Servers / 5. Fileservers
	// Use the http.NewServeMux's .Handle() method to add a handler

	// 2. Routing / 1. Stateful Handlers
	// Wrap the http.FileServer handler with the middleware "middlewareMetricsInc"
	// mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	app_router.Handle("/app", fsHandler)
	app_router.Handle("/app/*", fsHandler)

	// I recommend using the mux.HandleFunc to register your handler.
	// mux.HandleFunc("/healthz", handlerReadiness)
	app_router.Get("/healthz", handlerReadiness)

	app_router.Get("/metrics", apiCfg.handlerMetrics)
	// register a handler on the /reset path
	app_router.Get("/reset", apiCfg.handlerReset)

	// 1. Servers / 4. Server
	//mux := http.NewServeMux()
	api_router := chi.NewRouter()

	// 1. Servers / 11. Custom Handlers
	// The endpoint should be accessible at the /healthz path using any HTTP method.
	// api_router.Get("/healthz", handlerReadiness)
	// api_router.Get("/reset", apiCfg.handlerReset)

	api_router.Post("/revoke", apiCfg.handlerRevoke)
	api_router.Post("/refresh", apiCfg.handlerRefresh)
	api_router.Post("/login", apiCfg.handlerLogin)
	api_router.Post("/users", apiCfg.handlerUsersCreate)

	api_router.Put("/users", apiCfg.handlerUsersUpdate)

	api_router.Post("/chirps", apiCfg.handlerChirpsCreate)

	api_router.Get("/chirps", apiCfg.handlerChirpsRetrieve)
	api_router.Get("/chirps/{chirpID}", apiCfg.handlerChirpsGet)

	api_router.Delete("/chirps/{chirpID}", apiCfg.handlerChirpsDelete)

	api_router.Post("/polka/webhooks", apiCfg.handlerWebhook)

	app_router.Mount("/api", api_router)

	// 1. Servers / 4. Server
	// mux := http.NewServeMux()
	admin_router := chi.NewRouter()

	// admin_router.Get("/metrics", apiCfg.handlerMetrics)

	app_router.Mount("/admin", admin_router)

	// 1. Servers / 4. Server
	// Wrap that mux in a custom middleware function that adds CORS headers
	// corsMux := middlewareCors(mux)
	corsMux := middlewareCors(app_router)

	// 1. Servers / 4. Server
	// Create a new http.Server and use the corsMux as the handler
	// srv := &http.Server{
	// 	Addr:    ":" + port,
	// 	Handler: corsMux,
	// }
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	// 1. Servers / 4. Server
	// Use the server's ListenAndServe method to start the server
	// log.Printf("Serving on port: %s\n", port)
	// log.Fatal(srv.ListenAndServe())
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
