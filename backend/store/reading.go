package store

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nick-moyer/seed-sentinel/models"
)

func CalculateMoisturePercentage(rawValue, dryRef, wetRef int) int {
	if rawValue <= dryRef {
		return 0
	}

	if rawValue >= wetRef {
		return 100
	}

	return (wetRef - rawValue) * 100 / (wetRef - dryRef)
}

// Saves a sensor reading
func InsertReading(data models.SensorReadingPayload) {
	// Fetch calibration values
	dryRef, wetRef, err := GetCalibration(data.SensorID)
	if err != nil {
		log.Printf("Sensor with id %s does not exist or calibration not found. Reading not saved.\n", data.SensorID)
		return
	}

	// Fetch plant by sensor ID
	plant, err := FetchPlantBySensorID(data.SensorID)
	if err != nil {
		log.Printf("No plant found for sensor %s. Reading not saved.\n", data.SensorID)
		return
	}

	// Convert raw value to moisture percentage
	moisture := CalculateMoisturePercentage(data.RawValue, dryRef, wetRef)

	// Insert reading into DB
	stmt, err := db.Prepare("INSERT INTO readings(plant_id, moisture_percentage) VALUES(?, ?)")
	if err != nil {
		log.Println("Database Error (Prepare):", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(plant.ID, moisture)
	if err != nil {
		log.Println("Database Error (Insert):", err)
	} else {
		log.Println("Saved reading to DB")
	}
}

// Returns up to `limit` readings for a sensor
func FetchReadings(sensorID string, limit int) ([]models.Reading, error) {
	if limit < 1 {
		limit = 100
	}

	query := fmt.Sprintf(`
        SELECT r.moisture_percentage, r.created_at
        FROM readings r
        INNER JOIN plants p ON r.plant_id = p.id
        WHERE p.sensor_id = ?
        ORDER BY r.created_at ASC
        LIMIT %d
    `, limit)

	rows, err := db.Query(query, sensorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.Reading
	for rows.Next() {
		var reading models.Reading
		if err := rows.Scan(&reading.MoisturePercentage, &reading.CreatedAt); err != nil {
			return nil, err
		}
		history = append(history, reading)
	}
	if history == nil {
		history = []models.Reading{}
	}
	return history, nil
}
