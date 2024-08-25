package httputil

// ==== httputils: requests.go =====================================================================================
// For inspecting + debugging http.Requests into a readable structure & format
// ====================================================================================================================

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// RequestDetails is a struct to hold details about an http.Request
type RequestDetails struct {
	Method     string            `json:"method,omitempty"`
	Path       string            `json:"path,omitempty"`
	RemoteAddr string            `json:"remoteAddr,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Query      map[string]string `json:"query,omitempty"`
	Body       string            `json:"body,omitempty"`
	Timestamp  string            `json:"timestamp,omitempty"`
}

// Create a RequestDetails struct from an http.Request
func NewRequestDetails(r *http.Request, readBody bool) RequestDetails {
	headers := make(map[string]string)
	for k, v := range r.Header {
		headers[k] = strings.Join(v, ",")
	}

	query := make(map[string]string)
	for k, v := range r.URL.Query() {
		query[k] = strings.Join(v, ",")
	}

	bodyStr := ""

	if readBody {
		// Read the body if bodyDebug is enabled
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}

		bodyStr = string(body)

		// Reset the body so it can be read again!
		r.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	return RequestDetails{
		Method:     r.Method,
		Path:       r.URL.Path,
		RemoteAddr: r.RemoteAddr,
		Headers:    headers,
		Query:      query,
		Body:       bodyStr,
		Timestamp:  time.Now().Format(time.RFC3339),
	}
}
