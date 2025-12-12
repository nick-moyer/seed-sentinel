package models

import "time"

// Data received from sensor
type SensorConfigPayload struct {
	SensorID  string `json:"sensor_id"`
	PlantName string `json:"plant_name"`
}

type SensorReadingPayload struct {
	SensorID  string `json:"sensor_id"`
	Moisture  int    `json:"moisture"`
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
