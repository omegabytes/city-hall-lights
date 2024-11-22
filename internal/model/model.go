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

type ImageMetadata struct {
	FileName    string      `json:"file_name"`
	AltText     string      `json:"alt_text"`
	Attribution Attribution `json:"attribution"`
}

type Attribution struct {
	Creator    string `json:"creator"`
	Title      string `json:"title"`
	SourceURL  string `json:"source_url"`
	LicenseURL string `json:"license_url"`
}
