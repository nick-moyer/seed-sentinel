package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// 1. Define the Data Contract (What the ESP32 will send)
type Telemetry struct {
    SensorID string  `json:"sensor_id"`
    Moisture int     `json:"moisture"`
}

func main() {
    // 2. Define the Handler
    http.HandleFunc("/telemetry", func(w http.ResponseWriter, r *http.Request) {
        // Only allow POST
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        // Decode JSON
        var data Telemetry
        if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
            http.Error(w, "Bad JSON", http.StatusBadRequest)
            return
        }

        // Log it (This proves it works!)
        fmt.Printf("[%s] Received: Plant=%s Moisture=%d%%\n",
            time.Now().Format(time.RFC3339), data.SensorID, data.Moisture)

        // Send 200 OK back to the sensor
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ACK"))
    })

    // 3. Start Server
    fmt.Println("🌱 Seed Sentinel Backend listening on :8080...")
    http.ListenAndServe(":8080", nil)
}