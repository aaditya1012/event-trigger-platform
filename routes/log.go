package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"event-trigger-platform/db"
)

func ListLogs(w http.ResponseWriter, r *http.Request) {
	// Fetch logs from the database
	logs, err := db.GetEventLogs()
	if err != nil {
		log.Printf("Error fetching event logs: %v", err)
		http.Error(w, "Failed to fetch logs", http.StatusInternalServerError)
		return
	}

	// Respond with the logs
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
