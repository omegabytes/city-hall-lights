package model

import "time"

type Events struct {
	Events []Event `json:"events"`
}

type Event struct {
	DateString     string    `json:"date_string"`
	StartTimeStamp time.Time `json:"start_timestamp"`
	Color          string    `json:"color"`
	Description    string    `json:"description"`
	RawEventString string    `json:"raw_event_string"`
}
