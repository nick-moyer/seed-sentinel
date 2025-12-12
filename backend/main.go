package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/nick-moyer/seed-sentinel/models"
	"github.com/nick-moyer/seed-sentinel/services"
	"github.com/nick-moyer/seed-sentinel/store"
)

// Helper to enforce POST and decode JSON
func handlePostJSON[T any](w http.ResponseWriter, r *http.Request, payload *T) bool {
	// Only allow POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}

	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return false
	}
	return true
}

// --- MAIN ---

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Warning: .env file not found")
	}

	// Initialize DB
	store.InitDB()

	// HTTP Handlers
	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		// Parse request
		var data models.SensorConfigPayload
		if !handlePostJSON(w, r, &data) {
			return
		}

		// Log it
		fmt.Printf("[%s] Received Config: Sensor=%s Plant=%s\n",
			time.Now().Format(time.RFC3339), data.SensorID, data.PlantName)

		// Save to DB
		store.SaveSensor(data.SensorID, data.PlantName)

		// Send 200 OK back to the sensor
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ACK"))
	})

	http.HandleFunc("/telemetry", func(w http.ResponseWriter, r *http.Request) {
		// Parse request
		var data models.SensorReadingPayload
		if !handlePostJSON(w, r, &data) {
			return
		}

		// Log it
		fmt.Printf("[%s] Received: Sensor=%s Moisture=%d%%\n",
			time.Now().Format(time.RFC3339), data.SensorID, data.Moisture)

		// Save to DB
		store.SaveReading(data)

		// Run LLM-Agent in background
		go func(data models.SensorReadingPayload) {
			decision, _ := services.RunAgent(data)
			if decision.AlertNeeded {
				services.SendNotification(decision.Advice)
			}
		}(data)

		// Send 200 OK back to the sensor
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ACK"))
	})

	// Start Server
	fmt.Println("Seed Sentinel Backend listening on :8080...")
	http.ListenAndServe(":8080", nil)
}
