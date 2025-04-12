package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"smartraveller-api/smartraveller"
	"time"
)

type AdvisoryResponse struct {
	Fetched  string                 `json:"lastFetched"`
	Advisory smartraveller.Advisory `json:"advisory"`
}

func GetAdvisory(w http.ResponseWriter, r *http.Request) {
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

	// # https://vercel.com/docs/edge-network/caching
	w.Header().Set("Cache-Control", "s-maxage=1800")
	w.Header().Set("Content-Type", "application/json")

	response := AdvisoryResponse{
		Fetched:  time.Now().UTC().Format(time.RFC3339),
		Advisory: advisories[0],
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
