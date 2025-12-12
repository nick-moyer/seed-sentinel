package models

import "time"

// Data received from sensor
type SensorPayload struct {
	SensorID  string `json:"sensor_id"`
	Moisture  int    `json:"moisture"`
	PlantName string `json:"plant_name"`
}

// Data stored in DB
type Reading struct {
	ID        int
	SensorID  int
	Moisture  int
	Timestamp time.Time
}

// Response from LLM-Agent
type AgentResponse struct {
	AlertNeeded bool   `json:"alert_needed"`
	Advice      string `json:"advice"`
}
