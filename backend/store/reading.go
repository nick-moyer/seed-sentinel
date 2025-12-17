package store

import (
	"context"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nick-moyer/seed-sentinel/models"
)

// Converts raw sensor value to moisture percentage using calibration data
func CalculateMoisturePercentage(rawValue, dryRef, wetRef int) int {
	if dryRef > wetRef {
		// Typical: higher value = drier
		if rawValue >= dryRef {
			return 0
		}
		if rawValue <= wetRef {
			return 100
		}
		return (dryRef - rawValue) * 100 / (dryRef - wetRef)
	} else {
		// Inverted
		if rawValue <= dryRef {
			return 0
		}
		if rawValue >= wetRef {
			return 100
		}
		return (rawValue - dryRef) * 100 / (wetRef - dryRef)
	}
}

// Returns up to `limit` readings for a sensor
func FetchReadings(ctx context.Context, sensorID string, limit int) ([]models.Reading, error) {
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

	rows, err := db.QueryContext(ctx, query, sensorID)
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

// Inserts a new sensor reading and returns the calculated moisture percentage
func InsertReading(ctx context.Context, data models.SensorReadingPayload) (int, error) {
	// Fetch calibration values
	dryRef, wetRef, err := FetchSensorCalibration(ctx, data.SensorID)
	if err != nil {
		return 0, fmt.Errorf("sensor with id %s does not exist or calibration not found: %w", data.SensorID, err)
	}

	// Fetch plant by sensor ID
	plant, err := FetchPlantBySensorID(ctx, data.SensorID)
	if err != nil {
		return 0, fmt.Errorf("no plant found for sensor %s: %w", data.SensorID, err)
	}

	// Convert raw value to moisture percentage
	moisture := CalculateMoisturePercentage(data.RawValue, dryRef, wetRef)

	// Insert reading into DB
	stmt, err := db.PrepareContext(ctx, "INSERT INTO readings(plant_id, moisture_percentage) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("database error (prepare): %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, plant.ID, moisture)
	if err != nil {
		return 0, fmt.Errorf("database error (insert): %w", err)
	}

	return moisture, nil
}
