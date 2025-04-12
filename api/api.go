package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"smartraveller-api/smartraveller"
	"time"
)

type AdvisoryResponse struct {
	Fetched  string                 `json:"lastFetched"`
	Advisory smartraveller.Advisory `json:"advisory"`
}

type AdvisoriesResponse struct {
	Fetched    string                   `json:"lastFetched"`
	Advisories []smartraveller.Advisory `json:"advisories"`
	Length     int                      `json:"length"`
}

func getIndex(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "<h1>Welcome to the Smartraveller API</h1>")
}

func getAdvisory(w http.ResponseWriter, r *http.Request) {
	countryAlpha2 := r.URL.Query().Get("country")

	if countryAlpha2 == "" || len(countryAlpha2) > 2 {
		http.Error(w, "Country code is invalid. Expected ISO 3166-1 alpha-2 country code.", http.StatusBadRequest)
		return
	}

	advisories, err := smartraveller.GetAdvisories(countryAlpha2)

	if err != nil {
		if errors.Is(err, smartraveller.ErrAdvisoryNotFound) {
			http.NotFound(w, r)
			return
		}

		http.Error(w, "Failed to fetch advisory", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := AdvisoryResponse{
		// get current time in UTC
		Fetched:  time.Now().UTC().Format(time.RFC3339),
		Advisory: advisories[0],
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func getAdvisories(w http.ResponseWriter, r *http.Request) {
	advisories, err := smartraveller.GetAdvisories("")

	if err != nil {
		http.Error(w, "Failed to fetch advisories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := AdvisoriesResponse{
		// get current time in UTC
		Fetched:    time.Now().UTC().Format(time.RFC3339),
		Advisories: advisories,
		Length:     len(advisories),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", getIndex)
	http.HandleFunc("/advisory", getAdvisory)
	http.HandleFunc("/advisories", getAdvisories)

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Server started successfully")
}
