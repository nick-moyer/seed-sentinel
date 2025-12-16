package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nick-moyer/seed-sentinel/models"
)

func RunAgent(payload models.AgentPayload) (models.AgentResponse, error) {
	// Prepare data for LLM-Agent
	data := map[string]any{
		"plant_name":          payload.PlantName,
		"plant_age_days":      payload.PlantAgeDays,
		"moisture_percentage": payload.MoisturePercentage,
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
