package main

import "net/http"

// 2. Routing / 1. Stateful Handlers
// Create a new handler that, when hit, will reset your fileserverHits back to 0
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
