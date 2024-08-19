package main

// ==== http-toolkit: static.go =======================================================================================
// Handlers for the purposes of serving static files and SPAs
// ====================================================================================================================

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var indexFile = "index.html"

// spaServe will serve content as static files from the configured directory
// It will fallback to serving the index file if the requested file does not exist
// This is to support SPAs that use client-side routing
func spaServe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")

	// Remove the route prefix from the URL path
	routePrefixNoSlash := strings.TrimSuffix(cfg.routePrefix, "/")
	r.URL.Path = strings.ReplaceAll(r.URL.Path, routePrefixNoSlash, "")

	// Get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// If we failed to get the absolute path respond with a 400 bad request and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Prepend the path with the path to the static directory
	path = filepath.Join(cfg.spaPath, path)

	// Check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// File does not exist, serve our index file
		http.ServeFile(w, r, filepath.Join(cfg.spaPath, indexFile))
		return
	} else if err != nil {
		// If we got a different error, something went super wrong
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Otherwise, use http.FileServer to serve the file
	http.FileServer(http.Dir(cfg.spaPath)).ServeHTTP(w, r)
}

// staticServe will serve content as static files from the configured directory
// It will fallback to serving directories listing if the path is a directory
func staticServe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")

	// Remove the route prefix from the URL path
	routePrefixNoSlash := strings.TrimSuffix(cfg.routePrefix, "/")
	r.URL.Path = strings.ReplaceAll(r.URL.Path, routePrefixNoSlash, "")

	http.FileServer(http.Dir(cfg.staticPath)).ServeHTTP(w, r)
}
