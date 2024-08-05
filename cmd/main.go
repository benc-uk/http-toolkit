package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

var cfg AppConfig
var tokenAuth *jwtauth.JWTAuth

func main() {
	// Set up configuration
	cfg = NewConfig()
	cfg.loadEnv()

	r := chi.NewRouter()

	if cfg.reqDebug {
		r.Use(reqDebugMiddleware)
	}

	r.Use(middleware.Logger)

	// Add all our routes under a sub-router
	// This allows us to set a custom prefix for all routes
	r.Route(cfg.routePrefix, func(r chi.Router) {
		r.Get("/", ok)
		r.Get("/health*", ok)

		r.Get("/info", systemInfo)

		r.HandleFunc("/status/{code}", statusCode)
		r.Get("/word", randomWord)
		r.Get("/word/{count}", randomWord)
		r.Get("/number", randomNumber)
		r.Get("/number/{max}", randomNumber)
		r.Get("/uuid", randomUUID)
		r.Get("/uuid/{input}", randomUUID)

		// Route protected by basic auth
		r.Route("/auth/basic", func(r chi.Router) {
			r.Use(middleware.BasicAuth("realm", map[string]string{
				cfg.basicAuthUser: cfg.basicAuthPassword,
			}))

			r.HandleFunc("/", ok)
		})

		// Route protected by simple SHA256 JWT auth
		r.Route("/auth/jwt", func(r chi.Router) {
			tokenAuth = jwtauth.New("HS256", []byte(cfg.jwtSignKey), nil)

			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))

			r.HandleFunc("/", ok)
		})

		// Handle fallback
		if cfg.inspectAll {
			// Add a catch-all route to inspect & echo requests that don't match any other routes
			r.HandleFunc("/*", inspect)
		} else {
			// Only inspect & echo requests to /inspect and /echo
			r.HandleFunc("/inspect", inspect)
			r.HandleFunc("/echo", inspect)
		}
	})

	server := &http.Server{
		Addr:              ":" + cfg.port,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       120 * time.Second,
		Handler:           r,
	}

	log.Printf("🌐 HTTP Toolkit v0.0")
	log.Printf("🚀 Server started on port %s", cfg.port)

	// Generate a valid JWT token for testing with no claims
	_, exampleToken, _ := tokenAuth.Encode(map[string]interface{}{})

	log.Printf("🔐 Basic auth credentials: %s:%s\n", cfg.basicAuthUser, cfg.basicAuthPassword)
	log.Printf("🔑 JWT valid token: %s\n\n", exampleToken)

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
