package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/joho/godotenv"
)

// --- STRUCTS ---

// What the ESP32 will send
type SensorPayload struct {
	SensorID  string `json:"sensor_id"`
	Moisture  int    `json:"moisture"`
	PlantName string `json:"plant_name"`
}

type AgentResponse struct {
	AlertNeeded bool   `json:"alert_needed"`
	Advice      string `json:"advice"`
}

// --- MAIN ---

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Warning: .env file not found")
	}

    // HTTP Handler
	http.HandleFunc("/telemetry", func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Decode JSON
		var data SensorPayload
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Bad JSON", http.StatusBadRequest)
			return
		}

		// Log it (This proves it works!)
		fmt.Printf("[%s] Received: Sensor=%s Moisture=%d%% Plant=%s\n",
			time.Now().Format(time.RFC3339), data.SensorID, data.Moisture, data.PlantName)

		// Send 200 OK back to the sensor
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ACK"))

		// Process Logic Asynchronously
		go processLogic(data)
	})

    // Start Server
	fmt.Println("🌱 Seed Sentinel Backend listening on :8080...")
	http.ListenAndServe(":8080", nil)
}

// --- HELPERS ---

func processLogic(data SensorPayload) {
	// 1. Send to Python
	jsonData, _ := json.Marshal(data)
	resp, err := http.Post("http://localhost:5000/analyze", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("⚠️  Agent is offline:", err)
		return
	}
	defer resp.Body.Close()

	// 2. Read Decision
	var agentResp AgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&agentResp); err != nil {
		fmt.Println("⚠️  Error decoding agent response:", err)
		return
	}

	// 3. Act
	if agentResp.AlertNeeded {
		fmt.Printf("🚨 ALERT: %s\n", agentResp.Advice)
		SendNotification(agentResp.Advice)
	} else {
		fmt.Println("✅ Status OK")
	}
}

func SendNotification(message string) {
	target := os.Getenv("NOTIFICATION_TARGET")

	if target == "" {
		fmt.Println("❌ Error: NOTIFICATION_TARGET is missing in .env")
		return
	}

	fmt.Printf("🔔 Sending alert to configured target...\n")

	url := fmt.Sprintf("https://ntfy.sh/%s", target)
	http.Post(url, "text/plain", strings.NewReader(message))
}
