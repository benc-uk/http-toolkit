package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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

	mux := http.NewServeMux()

	mux.HandleFunc("POST /inspect", inspect)
	mux.HandleFunc("GET /inspect", inspect)
	mux.HandleFunc("DELETE /inspect", inspect)
	mux.HandleFunc("PUT /inspect", inspect)
	mux.HandleFunc("PATCH /inspect", inspect)

	mux.HandleFunc("GET /info", systemInfo)
	mux.HandleFunc("GET /status/{code}", statusCode)

	mux.HandleFunc("GET /", ok)
	mux.HandleFunc("GET /health", ok)

	mux.HandleFunc("GET /word", randomWord)
	mux.HandleFunc("GET /word/{count}", randomWord)
	mux.HandleFunc("GET /number/{max}", randomNumber)
	mux.HandleFunc("GET /number", randomNumber)

	// Wrap the mux with the logger middleware
	loggerMux := NewLogger(mux)

	server := &http.Server{
		Addr:              ":" + portEnv,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           loggerMux,
	}

	log.Println("HTTP Debugger started on port " + portEnv)
	log.Fatal(server.ListenAndServe())
}

// Logger is a middleware handler that does request logging
type Logger struct {
	handler http.Handler
}

// ServeHTTP handles the request by passing it to the real
// handler and logging the request details
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))

	if reqDebug {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(newRequestDetails(r))
	}

	l.handler.ServeHTTP(w, r)
}

// NewLogger constructs a new Logger middleware handler
func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{handlerToWrap}
}
