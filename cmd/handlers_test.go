package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
