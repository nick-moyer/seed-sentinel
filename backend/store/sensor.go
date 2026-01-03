package store

import (
	"context"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nick-moyer/seed-sentinel/models"
)

func FetchAllSensors(ctx context.Context) ([]models.Sensor, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, dry_reference, wet_reference, created_at FROM sensors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensors []models.Sensor
	for rows.Next() {
		var sensor models.Sensor
		if err := rows.Scan(&sensor.ID, &sensor.DryReference, &sensor.WetReference, &sensor.CreatedAt); err != nil {
			return nil, err
		}
		sensors = append(sensors, sensor)
	}
	return sensors, nil
}

// Fetches dry and wet calibration values for a sensor by ID
func FetchSensorCalibration(ctx context.Context, sensorID string) (int, int, error) {
	var dryRef, wetRef int
	err := db.QueryRowContext(ctx, "SELECT dry_reference, wet_reference FROM sensors WHERE id = ?", sensorID).Scan(&dryRef, &wetRef)
	if err != nil {
		return 0, 0, err
	}
	return dryRef, wetRef, nil
}

// Updates or inserts a sensor calibration
func UpsertSensor(ctx context.Context, data models.SensorCalibrationPayload) error {
	stmt, err := db.PrepareContext(ctx, `
        INSERT INTO sensors (id, dry_reference, wet_reference)
        VALUES (?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            dry_reference = excluded.dry_reference,
            wet_reference = excluded.wet_reference
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, data.SensorID, data.DryReference, data.WetReference)
	return err
}
