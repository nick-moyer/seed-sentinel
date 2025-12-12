package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nick-moyer/seed-sentinel/models"
	"github.com/nick-moyer/seed-sentinel/store"
)

func RunAgent(payload models.SensorReadingPayload) (models.AgentResponse, error) {
	 // Query plant name from sensors table
    var plantName string
    err := store.DB().QueryRow("SELECT plant_name FROM sensors WHERE id = ?", payload.SensorID).Scan(&plantName)
    if err != nil {
        fmt.Println("Error fetching plant name:", err)
        plantName = "Unknown" // fallback if not found
    }

	// Prepare data for LLM-Agent
	data := map[string]any{
		"sensor_id": payload.SensorID,
		"moisture":  payload.Moisture,
		"plant_name": plantName,
	}

	jsonData, _ := json.Marshal(data)
	resp, err := http.Post("http://localhost:5000/analyze", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Agent is offline:", err)
		return models.AgentResponse{}, err
	}
	defer resp.Body.Close()

	var agentResp models.AgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&agentResp); err != nil {
		fmt.Println("Error decoding agent response:", err)
		return models.AgentResponse{}, err
	}
	return agentResp, nil
}
