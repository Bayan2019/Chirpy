package main

import "net/http"

// 1. Servers / 11. Custom Handlers
// Let's add a readiness endpoint to the Chirpy server!

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	// The endpoint should return a Content-Type: text/plain; charset=utf-8 header
	// Write the Content-Type header
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	// The endpoint should simply return a 200 OK status code
	// Write the status code using w.WriteHeader
	w.WriteHeader(http.StatusOK)
	// the body will contain a message that simply says "OK"
	// Write the body text using w.Write
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
