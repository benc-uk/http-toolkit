// Created by Copilot, don't blame me if the code is shonky!

package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func TestOkHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/ok", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ok)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestSystemInfoHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/system-info", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(systemInfo)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	expectedContentType := "application/json"
	if ct := rr.Header().Get("Content-Type"); ct != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v, want %v", ct, expectedContentType)
	}

	// Check that response has uptime, memory usage, and other system info
	body := rr.Body.String()
	if len(body) == 0 {
		t.Error("handler returned empty response body")
	}

	// Check that response body is valid JSON
	var sysInfo SystemInfo
	if err := json.Unmarshal([]byte(body), &sysInfo); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	// Check that response body has uptime, hostname, and memory usage
	if sysInfo.Uptime == "" {
		t.Error("handler returned empty uptime")
	}
	if sysInfo.Hostname == "" {
		t.Error("handler returned empty hostname")
	}
	if sysInfo.Memory == "" {
		t.Error("handler returned empty memory")
	}
}
func TestDelayHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/delay/2", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	// Only way to add URL parameters to a request is to use chi's RouteContext
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("seconds", "2")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(delay)

	startTime := time.Now()
	handler.ServeHTTP(rr, req)
	endTime := time.Now()

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v, want %v", rr.Body.String(), expected)
	}

	duration := endTime.Sub(startTime)
	expectedDuration := 2 * time.Second
	if duration < expectedDuration {
		t.Errorf("handler did not delay for expected duration: got %v, want at least %v", duration, expectedDuration)
	}
}
func TestRandomUUID(t *testing.T) {
	// Test case 1: No input
	req, err := http.NewRequest(http.MethodGet, "/uuid", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(randomUUID)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Test case 2: Input provided
	req, err = http.NewRequest(http.MethodGet, "/uuid", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(randomUUID)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check that the response body is a valid UUID
	body := rr.Body.String()
	if _, err := uuid.Parse(body); err != nil {
		t.Errorf("handler returned invalid UUID: %v", err)
	}
}
