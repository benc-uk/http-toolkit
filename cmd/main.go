package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var reqDebug bool
var bodyDebug bool

func main() {
	// Get PORT environment variable
	portEnv := os.Getenv("PORT")
	if portEnv == "" {
		portEnv = "8080"
	}

	reqDebug = true

	reqDebugEnv := strings.ToLower(os.Getenv("REQUEST_DEBUG"))
	if reqDebugEnv == "false" {
		reqDebug = false
	}

	bodyDebug = true

	bodyDebugEnv := strings.ToLower(os.Getenv("BODY_DEBUG"))
	if bodyDebugEnv == "false" {
		bodyDebug = false
	}

	r := chi.NewRouter()

	if reqDebug {
		r.Use(reqDebugMiddleware)
	}

	r.Use(middleware.Logger)

	r.Get("/", ok)
	r.Get("/health*", ok)

	r.Get("/info", systemInfo)

	r.Get("/status/{code}", statusCode)
	r.Get("/word", randomWord)
	r.Get("/word/{count}", randomWord)
	r.Get("/number", randomNumber)
	r.Get("/number/{max}", randomNumber)
	r.Get("/uuid", randomUUID)
	r.Get("/uuid/{input}", randomUUID)

	// Add a catch-all route to inspect & echo requests
	r.HandleFunc("/*", inspect)

	server := &http.Server{
		Addr:              ":" + portEnv,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           r,
	}

	log.Printf("HTTP Toolkit v0.0")
	log.Printf("Server started on port %s", portEnv)
	log.Fatal(server.ListenAndServe())
}

// Middleware to log 'deep' request details to the console
func reqDebugMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Debug requests to JSON string and log to console
		reqDetails := newRequestDetails(r)

		reqJSON, err := json.MarshalIndent(reqDetails, "", "  ")
		if err != nil {
			log.Println(err)
		}

		log.Println("Request details:", string(reqJSON))

		next.ServeHTTP(w, r)
	})
}
