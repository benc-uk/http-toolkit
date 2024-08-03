//
// Created by Copilot
//

package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewRequestDetails(t *testing.T) {
	tests := []struct {
		name      string
		method    string
		url       string
		headers   map[string]string
		body      string
		bodyDebug bool
		wantBody  string
	}{
		{
			name:      "Basic GET request",
			method:    http.MethodGet,
			url:       "/test",
			headers:   map[string]string{"Content-Type": "application/json"},
			body:      "",
			bodyDebug: false,
			wantBody:  "",
		},
		{
			name:      "POST request with body and headers",
			method:    http.MethodPost,
			url:       "/test?query1=value1&query2=value2",
			headers:   map[string]string{"Content-Type": "application/json"},
			body:      `{"key":"value"}`,
			bodyDebug: true,
			wantBody:  `{"key":"value"}`,
		},
		{
			name:      "Empty body with bodyDebug enabled",
			method:    http.MethodPost,
			url:       "/test",
			headers:   map[string]string{"Content-Type": "application/json"},
			body:      "",
			bodyDebug: true,
			wantBody:  "",
		},
		{
			name:      "BodyDebug disabled",
			method:    http.MethodPost,
			url:       "/test",
			headers:   map[string]string{"Content-Type": "application/json"},
			body:      `{"key":"value"}`,
			bodyDebug: false,
			wantBody:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyDebug = tt.bodyDebug
			req := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			rd := newRequestDetails(req)

			if rd.Method != tt.method {
				t.Errorf("expected method %s, got %s", tt.method, rd.Method)
			}
			if rd.Path != req.URL.Path {
				t.Errorf("expected path %s, got %s", req.URL.Path, rd.Path)
			}
			if rd.Body != tt.wantBody {
				t.Errorf("expected body %s, got %s", tt.wantBody, rd.Body)
			}
			for k, v := range tt.headers {
				if rd.Headers[k] != v {
					t.Errorf("expected header %s: %s, got %s", k, v, rd.Headers[k])
				}
			}
			for k, v := range req.URL.Query() {
				if rd.Query[k] != strings.Join(v, ",") {
					t.Errorf("expected query %s: %s, got %s", k, strings.Join(v, ","), rd.Query[k])
				}
			}
			if _, err := time.Parse(time.RFC3339, rd.Timestamp); err != nil {
				t.Errorf("expected valid timestamp, got %s", rd.Timestamp)
			}
		})
	}
}
