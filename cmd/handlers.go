package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"math/rand"

	"github.com/elastic/go-sysinfo"
)

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
	_ = enc.Encode(newRequestDetails(r))
}

func ok(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/health" {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Not Found"))

		return
	}

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

	log.Println("Status code:", code)

	status, err := strconv.Atoi(code)
	if err != nil {
		status = http.StatusBadRequest
	}

	w.WriteHeader(status)
	_, _ = w.Write([]byte(http.StatusText(status)))
}

func randomWord(w http.ResponseWriter, r *http.Request) {
	count := r.PathValue("count")
	if count == "" {
		count = "1"
	}

	countInt, err := strconv.Atoi(count)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid count value"))

		return
	}

	// Generate a random words and append them to a string
	wordsOut := ""
	for i := 0; i < countInt; i++ {
		//nolint:gosec
		wordsOut += words[rand.Intn(len(words))] + " "
	}

	_, _ = w.Write([]byte(wordsOut))
}

func randomNumber(w http.ResponseWriter, r *http.Request) {
	max := r.PathValue("max")
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
