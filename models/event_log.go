package models

import "time"

type EventLog struct {
	TriggerID string    `json:"trigger_id"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

func (e EventLog) Printf(s string, d string) {
	panic("unimplemented")
}
