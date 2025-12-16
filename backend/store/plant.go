package store

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/nick-moyer/seed-sentinel/models"
)

// FetchPlantBySensorID returns the most recent plant for a given sensor_id
func FetchPlantBySensorID(sensorID string) (*models.Plant, error) {
	var plant models.Plant
	err := db.QueryRow(`SELECT id, sensor_id, name, date_planted, created_at, updated_at FROM plants WHERE sensor_id = ? ORDER BY id DESC LIMIT 1`, sensorID).
		Scan(&plant.ID, &plant.SensorID, &plant.Name, &plant.DatePlanted, &plant.CreatedAt, &plant.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &plant, nil
}
