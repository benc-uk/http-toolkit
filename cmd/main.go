package main

import (
	"crypto/tls"
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
var version = "0.0"

func main() {
	log.Printf("🌐 HTTP Toolkit " + version)

	// Set up configuration
	cfg = NewConfig()
	cfg.loadFlags()
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

		r.HandleFunc("/delay/{seconds}", delay)
		r.HandleFunc("/delay", delay)

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
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		Handler:           r,
	}

	// Generate a valid JWT token for testing with no claims
	_, exampleToken, _ := tokenAuth.Encode(map[string]interface{}{})

	log.Printf("🔐 Basic auth credentials: %s:%s\n", cfg.basicAuthUser, cfg.basicAuthPassword)
	log.Printf("🔑 JWT valid token: %s\n\n", exampleToken)

	// Start the server using TLS if configured
	if cfg.useTLS {
		log.Printf("🚀 Server started with TLS on port %s", cfg.port)

		server.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}

		// ListenAndServeTLS blocks so nothing after this will run
		log.Fatal(server.ListenAndServeTLS(cfg.certPath+"/cert.pem", cfg.certPath+"/key.pem"))
	}

	// Otherwise start the server without TLS
	log.Printf("🚀 Server started on port %s", cfg.port)
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
