package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/nick-moyer/seed-sentinel/models"
	"github.com/nick-moyer/seed-sentinel/services"
	"github.com/nick-moyer/seed-sentinel/store"
)

// --- HELPERS ---

func handleJSON[T any](w http.ResponseWriter, r *http.Request, payload *T) bool {
	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return false
	}
	return true
}

func handleDBError(w http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return true
	}
	return false
}

// --- HANDLERS ---

// POST /config
func configHandler(w http.ResponseWriter, r *http.Request) {
	var data models.SensorConfigPayload

	// Parse request
	if !handleJSON(w, r, &data) {
		return
	}

	// Log it
	log.Printf("[%s] Received Config: Sensor=%s Plant=%s\n",
		time.Now().Format(time.RFC3339), data.SensorID, data.PlantName)

	// Save to DB
	store.UpsertSensor(data.SensorID, data.PlantName)

	// Send 200 OK back to the sensor
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ACK"))
}

// POST /telemetry
func telemetryHandler(w http.ResponseWriter, r *http.Request) {
	var data models.SensorReadingPayload

	// Parse request
	if !handleJSON(w, r, &data) {
		return
	}

	// Log it
	log.Printf("[%s] Received: Sensor=%s Moisture=%d%%\n",
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
}

// GET /api/sensors
func listSensorsHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch from DB
	sensors, err := store.ListSensors()
	if handleDBError(w, err) {
		return
	}

	// Return JSON
	json.NewEncoder(w).Encode(sensors)
}

// GET /api/history/{id}
func historyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Fetch from DB
	history, err := store.FetchReadings(id)
	if handleDBError(w, err) {
		return
	}

	// Return JSON
	json.NewEncoder(w).Encode(history)
}

// --- ROUTES ---

func registerRoutes(r *mux.Router) {
	r.HandleFunc("/config", configHandler).Methods("POST")
	r.HandleFunc("/telemetry", telemetryHandler).Methods("POST")
	r.HandleFunc("/api/sensors", listSensorsHandler).Methods("GET")
	r.HandleFunc("/api/history/{id}", historyHandler).Methods("GET")
	fs := http.FileServer(http.Dir("../frontend"))
	r.PathPrefix("/").Handler(fs)
}

// --- MAIN ---

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize DB
	store.InitDB()

	// Setup Router and Routes
	r := mux.NewRouter()
	registerRoutes(r)

	// Start Server
	log.Println("Seed Sentinel Backend listening on :8080...")
	http.ListenAndServe(":8080", r)
}
