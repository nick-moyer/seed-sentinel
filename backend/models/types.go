package models

import "time"
// Data received from sensor
type SensorCalibrationPayload struct {
	SensorID     string `json:"sensor_id"`
	DryReference int    `json:"dry_reference"`
	WetReference int    `json:"wet_reference"`
}

type SensorReadingPayload struct {
	SensorID string `json:"sensor_id"`
	RawValue int    `json:"raw_value"`
}

// Data sent to sensor for configuration
type PlantConfigurationPayload struct {
	SensorID    string    `json:"sensor_id"`
	Name        string    `json:"name"`
	DatePlanted time.Time `json:"date_planted"`
}

// Data stored in DB
type Reading struct {
	ID                 int
	PlantID            int
	MoisturePercentage int
	CreatedAt          time.Time
}

type Plant struct {
	ID          int
	SensorID    string
	Name        string
	DatePlanted time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Payload sent to LLM-Agent
type AgentPayload struct {
	PlantName          string `json:"plant_name"`
	PlantAgeDays       int    `json:"plant_age_days"`
	MoisturePercentage int    `json:"moisture_percentage"`
}

// Response from LLM-Agent
type AgentResponse struct {
	AlertNeeded bool   `json:"alert_needed"`
	Advice      string `json:"advice"`
}
