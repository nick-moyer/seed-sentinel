package store

import (
	"context"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nick-moyer/seed-sentinel/models"
)

// Updates or inserts a plant configuration
func UpsertPlantConfiguration(ctx context.Context, data models.PlantConfigurationPayload) error {
	stmt, err := db.PrepareContext(ctx, `
        INSERT INTO plants (sensor_id, name, date_planted)
        VALUES (?, ?, ?)
        ON CONFLICT(sensor_id) DO UPDATE SET
            name = excluded.name,
            date_planted = excluded.date_planted
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, data.SensorID, data.Name, data.DatePlanted)
	return err
}

// FetchPlantBySensorID returns the most recent plant for a given sensor_id
func FetchPlantBySensorID(ctx context.Context, sensorID string) (*models.Plant, error) {
	var plant models.Plant
	err := db.QueryRowContext(ctx, `SELECT id, sensor_id, name, date_planted, created_at, updated_at FROM plants WHERE sensor_id = ? ORDER BY id DESC LIMIT 1`, sensorID).
		Scan(&plant.ID, &plant.SensorID, &plant.Name, &plant.DatePlanted, &plant.CreatedAt, &plant.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &plant, nil
}
