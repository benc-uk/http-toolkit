package main

// ==== http-toolkit: config.go =======================================================================================
// This file contains the configuration for the application along with support for
// command line flags and environment variables
// ====================================================================================================================

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	reqDebug          bool
	bodyDebug         bool
	inspectAll        bool
	routePrefix       string
	port              string
	basicAuthUser     string
	basicAuthPassword string
	jwtSignKey        string
	certPath          string
	useTLS            bool
	spaPath           string
	staticPath        string
}

// NewConfig creates a new AppConfig with all default values
func NewConfig() Config {
	return Config{
		reqDebug:          true,
		bodyDebug:         true,
		inspectAll:        true,
		routePrefix:       "/",
		port:              "8000",
		basicAuthUser:     "admin",
		basicAuthPassword: "secret",
		jwtSignKey:        "key_1234567890",
		certPath:          "",
		useTLS:            false,
		spaPath:           "",
		staticPath:        "",
	}
}

// Load the command line flags passed as arguments into the config
func (cfg *Config) loadFlags() {
	printVer := flag.Bool("version", false, "Print version and exit")

	flag.StringVar(&cfg.port, "port", cfg.port, "Port to listen on")
	flag.BoolVar(&cfg.reqDebug, "request-debug", cfg.reqDebug, "Enable request debugging/inspection in logs")
	flag.BoolVar(&cfg.bodyDebug, "body-debug", cfg.bodyDebug, "Include request body when debugging")
	flag.BoolVar(&cfg.inspectAll, "inspect-fallback", cfg.inspectAll, "Fallback to inspect & echo when no routes match")
	flag.StringVar(&cfg.routePrefix, "route-prefix", cfg.routePrefix, "Route prefix")
	flag.StringVar(&cfg.basicAuthUser, "basic-auth-user", cfg.basicAuthUser, "Basic auth username")
	flag.StringVar(&cfg.basicAuthPassword, "basic-auth-password", cfg.basicAuthPassword, "Basic auth password")
	flag.StringVar(&cfg.jwtSignKey, "jwt-sign-key", cfg.jwtSignKey, "Signing key for JWT")
	flag.StringVar(&cfg.certPath, "cert-path", cfg.certPath, "Path to directory with TLS cert & key files")
	flag.StringVar(&cfg.spaPath, "spa-path", cfg.spaPath,
		"Path to SPA files to serve, default is none and don't serve SPA")
	flag.StringVar(&cfg.staticPath, "static-path", cfg.staticPath,
		"Path to static files to serve, default is none and don't serve files")

	flag.Usage = func() {
		fmt.Printf("http-toolkit %s - A simple HTTP toolkit for debugging and testing", version)
		fmt.Println("Usage: http-toolkit [flags]")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *printVer {
		fmt.Println(version)
		os.Exit(0)
	}
}

// Load the environment variables into the config
//
//nolint:cyclop
func (cfg *Config) loadEnv() {
	port := os.Getenv("PORT")
	if port != "" {
		cfg.port = port
	}

	reqDebug := strings.ToLower(os.Getenv("REQUEST_DEBUG"))
	if reqDebug == "false" || reqDebug == "0" {
		cfg.reqDebug = false
	}

	bodyDebug := strings.ToLower(os.Getenv("BODY_DEBUG"))
	if bodyDebug == "false" || bodyDebug == "0" {
		cfg.bodyDebug = false
	}

	inspectAll := strings.ToLower(os.Getenv("INSPECT_FALLBACK"))
	if inspectAll == "false" || inspectAll == "0" {
		cfg.inspectAll = false
	}

	routePrefix := os.Getenv("ROUTE_PREFIX")
	if routePrefix != "" {
		cfg.routePrefix = routePrefix

		if !strings.HasSuffix(cfg.routePrefix, "/") {
			cfg.routePrefix += "/"
		}

		if !strings.HasPrefix(cfg.routePrefix, "/") {
			cfg.routePrefix = "/" + cfg.routePrefix
		}
	}

	basicAuthUser := os.Getenv("BASIC_AUTH_USER")
	if basicAuthUser != "" {
		cfg.basicAuthUser = basicAuthUser
	}

	basicAuthPassword := os.Getenv("BASIC_AUTH_PASSWORD")
	if basicAuthPassword != "" {
		cfg.basicAuthPassword = basicAuthPassword
	}

	jwtSignKey := os.Getenv("JWT_SIGN_KEY")
	if jwtSignKey != "" {
		cfg.jwtSignKey = jwtSignKey
	}

	certPath := os.Getenv("CERT_PATH")
	if certPath != "" {
		cfg.certPath = certPath
	}

	spaPath := os.Getenv("SPA_PATH")
	if spaPath != "" {
		cfg.spaPath = spaPath
	}

	staticPath := os.Getenv("STATIC_PATH")
	if staticPath != "" {
		cfg.staticPath = staticPath
	}

	cfg.useTLS = false

	// Check for TLS cert & key files if certPath is set
	if cfg.certPath != "" {
		log.Printf("🧬 Enabling TLS, checking cert & key files in: %s", cfg.certPath)
		cfg.useTLS = true

		// Check cert & key files exist
		if _, err := os.Stat(cfg.certPath + "/cert.pem"); os.IsNotExist(err) {
			log.Printf("😟 cert.pem not found in cert path, TLS will be disabled")

			cfg.useTLS = false
		}

		if _, err := os.Stat(cfg.certPath + "/key.pem"); os.IsNotExist(err) {
			log.Printf("😟 key.pem not found in cert path, TLS will be disabled")

			cfg.useTLS = false
		}
	}
}
