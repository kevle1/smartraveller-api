package api

import (
	"encoding/json"
	"net/http"
	"smartraveller-api/smartraveller"
	"time"
)

type AdvisoriesResponse struct {
	Fetched    string                   `json:"lastFetched"`
	Advisories []smartraveller.Advisory `json:"advisories"`
	Length     int                      `json:"length"`
}

func GetAdvisories(w http.ResponseWriter, r *http.Request) {
	advisories, err := smartraveller.GetAdvisories("")

	if err != nil {
		http.Error(w, "Failed to fetch advisories", http.StatusInternalServerError)
		return
	}

	// # https://vercel.com/docs/edge-network/caching
	w.Header().Set("Cache-Control", "s-maxage=3600")
	w.Header().Set("Content-Type", "application/json")

	response := AdvisoriesResponse{
		Fetched:    time.Now().UTC().Format(time.RFC3339),
		Advisories: advisories,
		Length:     len(advisories),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
