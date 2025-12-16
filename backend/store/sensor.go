package store

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nick-moyer/seed-sentinel/models"
)

// Fetches dry and wet calibration values for a sensor by ID
func GetCalibration(sensorID string) (int, int, error) {
	var dryRef, wetRef int
	err := db.QueryRow("SELECT dry_reference, wet_reference FROM sensors WHERE id = ?", sensorID).Scan(&dryRef, &wetRef)
	if err != nil {
		return 0, 0, err
	}
	return dryRef, wetRef, nil
}

// Updates or inserts a sensor configuration
func UpsertSensor(data models.SensorCalibrationPayload) {
	// Try to update first
	updateStmt, err := db.Prepare("UPDATE sensors SET dry_reference = ?, wet_reference = ? WHERE id = ?")
	if err != nil {
		log.Println("Database Error (Prepare Update):", err)
		return
	}
	defer updateStmt.Close()

	res, err := updateStmt.Exec(data.DryReference, data.WetReference, data.SensorID)
	if err != nil {
		log.Println("Database Error (Update):", err)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("Database Error (RowsAffected):", err)
		return
	}

	if rowsAffected == 0 {
		// No row updated, insert new
		insertStmt, err := db.Prepare("INSERT INTO sensors(id, dry_reference, wet_reference) VALUES(?, ?, ?)")
		if err != nil {
			log.Println("Database Error (Prepare Insert):", err)
			return
		}
		defer insertStmt.Close()

		_, err = insertStmt.Exec(data.SensorID, data.DryReference, data.WetReference)
		if err != nil {
			log.Println("Database Error (Insert):", err)
		} else {
			log.Println("Inserted new sensor to DB:", data.SensorID)
		}
	} else {
		log.Println("Updated sensor in DB:", data.SensorID)
	}
}

// Return all active sensors
func ListSensors() ([]map[string]string, error) {
	rows, err := db.Query("SELECT id, plant_name FROM sensors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensors []map[string]string
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		sensors = append(sensors, map[string]string{"hardware_id": id, "plant_name": name})
	}
	if sensors == nil {
		sensors = []map[string]string{}
	}
	return sensors, nil
}
