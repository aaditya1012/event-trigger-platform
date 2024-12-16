package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"event-trigger-platform/models"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB() error {
	var err error
	dsn := "root:Qwerty@123@tcp(localhost:3306)/eventdb"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %v", err)
	}

	log.Println("Database connection established")
	return nil
}

// GetDB provides the current database connection
func GetDB() *sql.DB {
	return db
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		db.Close()
		log.Println("Database connection closed")
	}
}

// SaveTrigger saves a new trigger to the database
func SaveTrigger(trigger models.Trigger) error {
	query := `
		INSERT INTO triggers (id, type, payload, scheduled_at, recurring, interval, test)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	payloadJSON, _ := json.Marshal(trigger.Payload) // Convert map to JSON
	_, err := db.Exec(
		query,
		trigger.ID,
		trigger.Type,
		payloadJSON,
		trigger.ScheduledAt,
		trigger.Recurring,
		trigger.Intervals.Seconds(), // Store interval in seconds
		trigger.Test,
	)
	if err != nil {
		return fmt.Errorf("error saving trigger: %v", err)
	}
	log.Printf("Trigger saved: %v", trigger.ID)
	return nil
}

// GetAllTriggers retrieves all triggers from the database
func GetAllTriggers() ([]models.Trigger, error) {
	query := `
		SELECT id, type, payload, scheduled_at, recurring, interval, ` + `test` + ` FROM triggers
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching triggers: %v", err)
	}
	defer rows.Close()

	var triggers []models.Trigger
	for rows.Next() {
		var trigger models.Trigger
		var payloadJSON string
		var intervalSeconds int64

		if err := rows.Scan(
			&trigger.ID,
			&trigger.Type,
			&payloadJSON,
			&trigger.ScheduledAt,
			&trigger.Recurring,
			&intervalSeconds,
			&trigger.Test,
		); err != nil {
			return nil, fmt.Errorf("error scanning trigger row: %v", err)
		}

		// Convert JSON payload to map
		if err := json.Unmarshal([]byte(payloadJSON), &trigger.Payload); err != nil {
			log.Printf("Error unmarshalling payload for trigger %v: %v", trigger.ID, err)
		}
		trigger.Intervals = time.Duration(intervalSeconds) * time.Second

		triggers = append(triggers, trigger)
	}
	return triggers, nil
}

// UpdateTrigger updates an existing trigger in the database
func UpdateTrigger(trigger models.Trigger) error {
	query := `
		UPDATE triggers
		SET scheduled_at = ?, payload = ?, recurring = ?, interval = ?, test = ?
		WHERE id = ?
	`
	payloadJSON, _ := json.Marshal(trigger.Payload)
	_, err := db.Exec(
		query,
		trigger.ScheduledAt,
		payloadJSON,
		trigger.Recurring,
		trigger.Intervals.Seconds(),
		trigger.Test,
		trigger.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating trigger: %v", err)
	}
	log.Printf("Trigger updated: %v", trigger.ID)
	return nil
}

// DeleteTrigger deletes a trigger from the database by ID
func DeleteTrigger(id string) error {
	query := `DELETE FROM triggers WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting trigger: %v", err)
	}
	log.Printf("Trigger deleted: %v", id)
	return nil
}

// SaveEventLog saves a new event log to the database
func SaveEventLog(log models.EventLog) error {
	query := `
		INSERT INTO event_logs (trigger_id, timestamp, message)
		VALUES (?, ?, ?)
	`
	_, err := db.Exec(query, log.TriggerID, log.Timestamp, log.Message)
	if err != nil {
		return fmt.Errorf("error saving event log: %v", err)
	}
	log.Printf("Event log saved for trigger: %v", log.TriggerID)
	return nil
}

// GetEventLogs retrieves all event logs from the database
func GetEventLogs() ([]models.EventLog, error) {
	query := `SELECT trigger_id, timestamp, message FROM event_logs`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching event logs: %v", err)
	}
	defer rows.Close()

	var logs []models.EventLog
	for rows.Next() {
		var log models.EventLog
		if err := rows.Scan(&log.TriggerID, &log.Timestamp, &log.Message); err != nil {
			return nil, fmt.Errorf("error scanning event log row: %v", err)
		}
		logs = append(logs, log)
	}
	return logs, nil
}
