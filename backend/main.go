package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/nick-moyer/seed-sentinel/models"
	"github.com/nick-moyer/seed-sentinel/services"
	"github.com/nick-moyer/seed-sentinel/store"
)

// --- HELPERS ---

type spaHandler struct {
	staticPath string
	indexPath  string
}

func handleJSON[T any](w http.ResponseWriter, r *http.Request, payload *T) bool {
	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return false
	}

	// Log it
	log.Printf("[%s] Received JSON: %+v\n", time.Now().Format(time.RFC3339), payload)

	return true
}

// --- HANDLERS ---

// GET /*
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get absolute path to prevent directory traversal
	path, err := filepath.Abs(h.staticPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Join the path with the requested URL
	path = filepath.Join(path, r.URL.Path)

	// Check if the file exists
	_, err = os.Stat(path)

	// MAGIC: If file doesn't exist (like /settings), serve index.html
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If file DOES exist (like main.js), serve it normally
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

// POST /calibrate
func calibrateHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.SensorCalibrationPayload
	ctx := r.Context()

	// Parse request
	if !handleJSON(w, r, &payload) {
		return
	}

	// Save to DB
	if err := store.UpsertSensor(ctx, payload); err != nil {
		http.Error(w, "Failed to save calibration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Calibrated sensor %s: Dry=%d, Wet=%d\n", payload.SensorID, payload.DryReference, payload.WetReference)

	// Send 200 OK back to the sensor
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ACK"))
}

// POST /configure
func configureHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.PlantConfigurationPayload
	ctx := r.Context()

	// Parse request
	if !handleJSON(w, r, &payload) {
		return
	}

	// Save to DB
	if err := store.UpsertPlantConfiguration(ctx, payload); err != nil {
		http.Error(w, "Failed to save plant configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Configured plant for sensor %s: %+v\n", payload.SensorID, payload.Name)

	// Send 200 OK back to the sensor
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ACK"))
}

// POST /telemetry
func telemetryHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.SensorReadingPayload
	ctx := r.Context()

	// Parse request
	if !handleJSON(w, r, &payload) {
		return
	}

	// Save to DB
	latest, err := store.InsertReading(ctx, payload)
	if err != nil {
		http.Error(w, "Failed to save reading: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch plant details
	plant, err := store.FetchPlantBySensorID(ctx, payload.SensorID)
	if err != nil {
		http.Error(w, "Failed to fetch plant details: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Log latest reading
	log.Printf("Latest reading for sensor %s: %+v\n", payload.SensorID, latest)

	// Run LLM-Agent in background
	go func(m models.AgentPayload) {
		decision, _ := services.RunAgent(m)
		if decision.AlertNeeded {
			services.SendNotification(decision.Advice)
		}
	}(models.AgentPayload{
		PlantName:          plant.Name,
		PlantAgeDays:       int(time.Since(plant.DatePlanted).Hours() / 24),
		MoisturePercentage: latest,
	})

	// Send 200 OK back to the sensor
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ACK"))
}

// --- ROUTES ---

func registerRoutes(r *mux.Router) {
	r.HandleFunc("/calibrate", calibrateHandler).Methods("POST")
	r.HandleFunc("/telemetry", telemetryHandler).Methods("POST")
	r.HandleFunc("/configure", configureHandler).Methods("POST")

	// We use PathPrefix("/") to catch EVERYTHING else
	spa := spaHandler{staticPath: "../frontend/build", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	http.ListenAndServe(":8080", r)
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
