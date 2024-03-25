package main

import (
	"fmt"
	"net/http"
)

// 2. Routing / 1. Stateful Handlers
// Create a new handler that writes the number of requests
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>
		<body>
			Hits: %d
		</body>
	</html>`, cfg.fileserverHits)))
}

// 2. Routing / 1. Stateful Handlers
// write a new middleware method on a *apiConfig that increments the fileserverHits counter
// // 2. Routing / 1. Middleware
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}
