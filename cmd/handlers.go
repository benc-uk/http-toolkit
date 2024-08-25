package main

// ==== http-toolkit: handlers.go =====================================================================================
// Contains all the HTTP handlers these are mounted on various routes in main.go
// ====================================================================================================================

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"math/rand"

	"github.com/benc-uk/http-toolkit/pkg/httputil"
	"github.com/benc-uk/http-toolkit/pkg/stringutil"
	"github.com/elastic/go-sysinfo"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Used by the info handler to generate some useful system information
type SystemInfo struct {
	Hostname     string `json:"hostname"`
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	CPUCount     int    `json:"cpuCount"`
	Memory       string `json:"memory"`
	GoVersion    string `json:"goVersion"`
	ClientAddr   string `json:"clientAddr"`
	ServerHost   string `json:"serverHost"`
	Uptime       string `json:"uptime"`
}

func inspect(w http.ResponseWriter, r *http.Request) {
	// Return a JSON response with the request details
	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(httputil.NewRequestDetails(r, cfg.bodyDebug))
}

func ok(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func systemInfo(w http.ResponseWriter, r *http.Request) {
	host, _ := sysinfo.Host()
	mem, err := host.Memory()
	memString := "Unknown"

	if err == nil {
		memString = fmt.Sprintf("%dGB", mem.Total/1000/1000/1000)
	}

	info := SystemInfo{
		Hostname:     host.Info().Hostname,
		GoVersion:    runtime.Version(),
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		CPUCount:     runtime.NumCPU(),
		Memory:       memString,
		ClientAddr:   r.RemoteAddr,
		ServerHost:   r.Host,
		Uptime:       host.Info().Uptime().String(),
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(info)
}

func statusCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		code = "200"
	}

	status, err := strconv.Atoi(code)
	if err != nil {
		status = http.StatusBadRequest
	}

	w.WriteHeader(status)
	_, _ = w.Write([]byte(http.StatusText(status)))
}

func randomWord(w http.ResponseWriter, r *http.Request) {
	count := chi.URLParam(r, "count")
	if count == "" {
		count = "1"
	}

	countInt, err := strconv.Atoi(count)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid count value"))

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(strings.Join(stringutil.RandomWords(countInt), " ")))
}

func randomNumber(w http.ResponseWriter, r *http.Request) {
	max := chi.URLParam(r, "max")
	if max == "" {
		max = "1000"
	}

	maxInt, err := strconv.Atoi(max)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid max value"))

		return
	}

	// Generate a random number between 0 and max
	w.WriteHeader(http.StatusOK)

	//nolint:sec
	_, _ = w.Write([]byte(strconv.Itoa(rand.Intn(maxInt))))
}

func randomUUID(w http.ResponseWriter, r *http.Request) {
	input := chi.URLParam(r, "input")

	var u uuid.UUID

	var err error

	if input != "" {
		// if input is less than 16 characters, pad it with '0's
		sourceString := input
		if len(input) < 16 {
			sourceString = input + strings.Repeat("0", 16-len(input))
		}

		reader := strings.NewReader(sourceString)

		u, err = uuid.NewRandomFromReader(reader)
		if err != nil {
			_, _ = w.Write([]byte("Error generating UUID"))

			return
		}
	} else {
		u, err = uuid.NewRandom()
		if err != nil {
			_, _ = w.Write([]byte("Error generating UUID"))

			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(u.String()))
}

func delay(w http.ResponseWriter, r *http.Request) {
	delay := chi.URLParam(r, "seconds")
	if delay == "" {
		delay = "1"
	}

	delayInt, err := strconv.Atoi(delay)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid delay value"))

		return
	}

	time.Sleep(time.Duration(delayInt) * time.Second)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
