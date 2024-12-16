package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"event-trigger-platform/db"
	"event-trigger-platform/models"
)

var (
	triggerMutex sync.Mutex
)

// CreateTrigger handles the creation of a new trigger
func CreateTrigger(w http.ResponseWriter, r *http.Request) {
	var newTrigger models.Trigger
	if err := json.NewDecoder(r.Body).Decode(&newTrigger); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newTrigger.ID = fmt.Sprintf("trigger-%d", time.Now().UnixNano())

	// Save trigger to the database
	if err := db.SaveTrigger(newTrigger); err != nil {
		log.Printf("Error saving trigger: %v", err)
		http.Error(w, "Failed to save trigger", http.StatusInternalServerError)
		return
	}

	log.Printf("Created new trigger: %+v", newTrigger)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTrigger)
}

// ListTriggers retrieves all triggers from the database
func ListTriggers(w http.ResponseWriter, r *http.Request) {
	triggers, err := db.GetAllTriggers()
	if err != nil {
		log.Printf("Error fetching triggers: %v", err)
		http.Error(w, "Failed to fetch triggers", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(triggers)
}

// DeleteTrigger deletes a trigger by ID
func DeleteTrigger(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Trigger ID is required", http.StatusBadRequest)
		return
	}

	// Delete trigger from the database
	if err := db.DeleteTrigger(id); err != nil {
		log.Printf("Error deleting trigger with ID %s: %v", id, err)
		http.Error(w, "Failed to delete trigger", http.StatusInternalServerError)
		return
	}

	log.Printf("Deleted trigger with ID: %s", id)
	w.WriteHeader(http.StatusOK)
}

// LogEvent saves an event log in the database
func LogEvent(triggerID string, message string) {
	eventLog := models.EventLog{
		TriggerID: triggerID,
		Timestamp: time.Now(),
		Message:   message,
	}

	if err := db.SaveEventLog(eventLog); err != nil {
		log.Printf("Error saving event log: %v", err)
	}
}

// StartEventProcessing continuously checks for scheduled triggers and processes them
func StartEventProcessing() {
	for {
		time.Sleep(1 * time.Second) // Check every second

		triggerMutex.Lock()
		triggers, err := db.GetAllTriggers()
		if err != nil {
			log.Printf("Error fetching triggers: %v", err)
			triggerMutex.Unlock()
			continue
		}

		for _, trigger := range triggers {
			if trigger.Type == models.Scheduled && time.Now().After(trigger.ScheduledAt) {
				log.Printf("Executing scheduled trigger: %s", trigger.ID)
				LogEvent(trigger.ID, "Scheduled trigger executed")

				if trigger.Recurring {
					// Reschedule the recurring trigger
					trigger.ScheduledAt = trigger.ScheduledAt.Add(time.Duration(trigger.Intervals) * time.Second)
					if err := db.UpdateTrigger(trigger); err != nil {
						log.Printf("Error updating recurring trigger: %v", err)
					}
				} else {
					// Remove non-recurring trigger
					if err := db.DeleteTrigger(trigger.ID); err != nil {
						log.Printf("Error deleting trigger: %v", err)
					}
				}
			}
		}
		triggerMutex.Unlock()
	}
}
