package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Bayan2019/chirpy/internal/database"
	"github.com/go-chi/chi/v5"

	// 6. Authentication / 6. Authentication with JWTs
	//Because we're using a file,
	//and not manually adding the variable to our session,
	//we'll need to use a library to load the environment variables.
	"github.com/joho/godotenv"
)

// 2. Routing / 1. Stateful Handlers
// Create a struct that will hold any stateful,
// in-memory data we'll need to keep track of
// we just need to keep track of the number of requests we've received (fileserverHits)
// 5. Storage / 1. Storage
// the most important part of your typical web application is
// the storage of data (DB)
// 6. Authentication / 6. Authentication with JWTs
// you'll need to create a secret for your server
// - the secret is used to sign and verify JWTs
// You'll want to store the jwtSecret in your apiConfig struct
// so that your handlers can access it (jwtSecret).
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

	// 6. Authentication / 6. Authentication with JWTs
	// Secrets shouldn't be stored in Git,
	// just in case anyone malicious gains access to your repository
	// Environment variables are simple key/value pairs
	// that are available to the programs you run.
	// we'll store the secret in a gitingore'd file called .env
	// by default, godotenv will look for a file named .env in the current directory
	godotenv.Load(".env")
	// 6. Authentication / 6. Authentication with JWTs
	// Then you can load the JWT_SECRET variable using the standard library like
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY environment variable is not set")
	}

	// 6. Authentication / 6. Authentication with JWTs
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	// 6. Authentication / 6. Authentication with JWTs
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

	// 1. Servers / 5. Fileservers
	// I recommend using the mux.HandleFunc to register your handler.
	// mux.HandleFunc("/healthz", handlerReadiness)
	// 2. Routing / 4. Routing
	// mux.HandleFunc("GET /healthz", handlerReadiness)
	// mux.HandleFunc("GET /api/healthz", handlerReadiness)
	// // app_router.Get("/healthz", handlerReadiness)

	// 2. Routing / 1. Stateful Handlers
	// mux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	// 2. Routing / 4. Routing
	// mux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
	// mux.HandleFunc("GET /api/metrics", apiCfg.handlerMetrics)
	// app_router.Get("/metrics", apiCfg.handlerMetrics)

	// register a handler on the /reset path
	// mux.HandleFunc("/reset", apiCfg.handlerReset)
	// 2. Routing / 4. Routing
	// mux.HandleFunc("GET /reset", apiCfg.handlerReset)
	// mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	// app_router.Get("/reset", apiCfg.handlerReset)

	// 1. Servers / 4. Server
	//mux := http.NewServeMux()
	api_router := chi.NewRouter()

	// 1. Servers / 11. Custom Handlers
	// The endpoint should be accessible at the /healthz path using any HTTP method.
	api_router.Get("/healthz", handlerReadiness)
	// 2. Routing / 1. Stateful Handlers
	// register a handler on the /reset path
	// mux.HandleFunc("/reset", apiCfg.handlerReset)
	api_router.Get("/reset", apiCfg.handlerReset)

	// api_router.Get("/metrics", apiCfg.handlerMetrics)

	// 4. JSON / 2. JSON
	// Add a new endpoint to the Chirpy API that accepts a POST request at /api/validate_chirp
	// mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	// 5. Storage 1. Storage
	// Delete the /api/validate_chirp endpoint
	// api_router.Post("/validate_chirp", apiCfg.handlerChirpsCreate)

	// api_router.Post("/revoke", apiCfg.handlerRevoke)
	api_router.Post("/refresh", apiCfg.handlerRefresh)
	// 6. Authentication / 1. Authentication with Passwords
	// create a new POST /api/login endpoint.
	// This endpoint should allow a user to login.
	// 6. Authentication / 6. Authentication with JWTs
	// Update the POST /api/login endpoint
	api_router.Post("/login", apiCfg.handlerLogin)
	// 5. Storage / 7. Users
	// Add a new endpoint to your server that allows users to be created.
	// 6. Authentication / 1. Authentication with Passwords
	// Update the POST /api/users endpoint
	api_router.Post("/users", apiCfg.handlerUsersCreate)

	// 6. Authentication / 6. Authentication with JWTs
	// create a new PUT /api/users endpoint
	// This endpoint should update a user's email and password
	// mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)
	api_router.Put("/users", apiCfg.handlerUsersUpdate)

	// 5. Storage / 1. Storage
	// This endpoint should accept a JSON payload with a body field.
	// If all goes well, respond with a 201 status code and the full chirp resource.
	api_router.Post("/chirps", apiCfg.handlerChirpsCreate)

	// 5. Storage / 1. Storage
	// This endpoint should return an array of all chirps in the file, ordered by id in ascending order.
	// mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	api_router.Get("/chirps", apiCfg.handlerChirpsRetrieve)
	// 5. Storage / 4. Get
	// Add a new endpoint to your server that allows users to get a single chirp by ID.
	// mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)
	api_router.Get("/chirps/{chirpID}", apiCfg.handlerChirpsGet)

	api_router.Delete("/chirps/{chirpID}", apiCfg.handlerChirpsDelete)

	api_router.Post("/polka/webhooks", apiCfg.handlerWebhook)

	app_router.Mount("/api", api_router)

	// 1. Servers / 4. Server
	// mux := http.NewServeMux()
	admin_router := chi.NewRouter()

	// 2. Routing / 1. Stateful Handlers
	// mux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	// mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	admin_router.Get("/metrics", apiCfg.handlerMetrics)

	app_router.Mount("/admin", admin_router)

	// 1. Servers / 4. Server
	// Wrap that mux in a custom middleware function that adds CORS headers
	// corsMux := middlewareCors(mux)
	corsMux := middlewareCors(app_router)

	// 1. Servers / 4. Server
	// Create a new http.Server and use the corsMux as the handler
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	// 1. Servers / 4. Server
	// Use the server's ListenAndServe method to start the server
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
