package store

import (
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nick-moyer/seed-sentinel/models"
)

// Save a sensor reading
func SaveReading(payload models.SensorReadingPayload) {
	// Check if the sensor exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM sensors WHERE id = ?)", payload.SensorID).Scan(&exists)
	if err != nil {
		log.Println("Database Error (Sensor Exists Check):", err)
		return
	}
	if !exists {
		log.Printf("Sensor with id %s does not exist. Reading not saved.\n", payload.SensorID)
		return
	}

	stmt, err := db.Prepare("INSERT INTO readings(sensor_id, moisture, timestamp) VALUES(?, ?, ?)")
	if err != nil {
		log.Println("Database Error (Prepare):", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(payload.SensorID, payload.Moisture, time.Now())
	if err != nil {
		log.Println("Database Error (Insert):", err)
	} else {
		log.Println("Saved reading to DB")
	}
}

// Returns up to 100 readings for a sensor
func FetchReadings(sensorID string) ([]models.Reading, error) {
	rows, err := db.Query("SELECT moisture, timestamp FROM readings WHERE sensor_id = ? ORDER BY timestamp ASC LIMIT 100", sensorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.Reading
	for rows.Next() {
		var reading models.Reading
		if err := rows.Scan(&reading.Moisture, &reading.Timestamp); err != nil {
			return nil, err
		}
		history = append(history, reading)
	}
	if history == nil {
		history = []models.Reading{}
	}
	return history, nil
}
