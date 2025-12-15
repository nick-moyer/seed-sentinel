package store

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Updates or inserts a sensor configuration
func UpsertSensor(id string, plantName string) {
	// Try to update first
	updateStmt, err := db.Prepare("UPDATE sensors SET plant_name = ? WHERE id = ?")
	if err != nil {
		log.Println("Database Error (Prepare Update):", err)
		return
	}
	defer updateStmt.Close()

	res, err := updateStmt.Exec(plantName, id)
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
		insertStmt, err := db.Prepare("INSERT INTO sensors(id, plant_name) VALUES(?, ?)")
		if err != nil {
			log.Println("Database Error (Prepare Insert):", err)
			return
		}
		defer insertStmt.Close()

		_, err = insertStmt.Exec(id, plantName)
		if err != nil {
			log.Println("Database Error (Insert):", err)
		} else {
			log.Println("Inserted new sensor to DB:", id)
		}
	} else {
		log.Println("Updated sensor in DB:", id)
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
