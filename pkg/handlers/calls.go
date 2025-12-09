package handlers

import (
	"SpeechAnalytics/pkg/repositories"
	"encoding/json"
	"log"
	"net/http"
)

func getCalls(w http.ResponseWriter, r *http.Request) {
	calls, err := repositories.GetAllCalls()

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(calls)
}
