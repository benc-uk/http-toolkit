package main

import (
	"flag"
	"os"
	"strings"
)

type AppConfig struct {
	reqDebug          bool
	bodyDebug         bool
	inspectAll        bool
	routePrefix       string
	port              string
	basicAuthUser     string
	basicAuthPassword string
	jwtSignKey        string
}

func NewConfig() AppConfig {
	return AppConfig{
		reqDebug:          true,
		bodyDebug:         true,
		inspectAll:        true,
		routePrefix:       "/",
		port:              "8080",
		basicAuthUser:     "admin",
		basicAuthPassword: "secret",
		jwtSignKey:        "key_1234567890",
	}
}

func (cfg *AppConfig) LoadFlags() {
	printVer := flag.Bool("version", false, "Print version and exit")

	flag.StringVar(&cfg.port, "port", cfg.port, "Port to listen on")
	flag.BoolVar(&cfg.reqDebug, "request-debug", cfg.reqDebug, "Enable request debugging/inspection in logs")
	flag.BoolVar(&cfg.bodyDebug, "body-debug", cfg.bodyDebug, "Include request body when debugging")
	flag.BoolVar(&cfg.inspectAll, "inspect-fallback", cfg.inspectAll, "Fallback to inspect & echo when no routes match")
	flag.StringVar(&cfg.routePrefix, "route-prefix", cfg.routePrefix, "Route prefix")
	flag.StringVar(&cfg.basicAuthUser, "basic-auth-user", cfg.basicAuthUser, "Basic auth username")
	flag.StringVar(&cfg.basicAuthPassword, "basic-auth-password", cfg.basicAuthPassword, "Basic auth password")
	flag.StringVar(&cfg.jwtSignKey, "jwt-sign-key", cfg.jwtSignKey, "Signing key for JWT")

	flag.Parse()

	if *printVer {
		println(version)
		os.Exit(0)
	}
}

//nolint:cyclop
func (cfg *AppConfig) loadEnv() {
	// Get PORT environment variable
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
}
