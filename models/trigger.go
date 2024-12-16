package models

import "time"

type TriggerType string

const (
	Scheduled TriggerType = "Scheduled"
	API       TriggerType = "API"
)

type Trigger struct {
	ID          string            `json:"id"`
	Type        TriggerType       `json:"type"`
	Payload     map[string]string `json:"payload,omitempty"`
	ScheduledAt time.Time         `json:"scheduled_at,omitempty"`
	Recurring   bool              `json:"recurring,omitempty"`
	Intervals    time.Duration     `json:"interval,omitempty"`
	Test        bool              `json:"test"`
}
