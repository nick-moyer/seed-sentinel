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

// --- MAIN ---

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Warning: .env file not found")
	}

	// Initialize DB
	store.InitDB()

	// HTTP Handler
	http.HandleFunc("/telemetry", func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Decode JSON
		var data models.SensorPayload
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Bad JSON", http.StatusBadRequest)
			return
		}

		// Log it
		fmt.Printf("[%s] Received: Sensor=%s Moisture=%d%% Plant=%s\n",
			time.Now().Format(time.RFC3339), data.SensorID, data.Moisture, data.PlantName)

		// Send 200 OK back to the sensor
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ACK"))

		go func(data models.SensorPayload) {
			decision, _ := services.RunAgent(data)
			if decision.AlertNeeded {
				services.SendNotification(decision.Advice)
			}
		}(data)

		// Save to DB
		store.SaveReading(data)
	})

	// Start Server
	fmt.Println("Seed Sentinel Backend listening on :8080...")
	http.ListenAndServe(":8080", nil)
}
