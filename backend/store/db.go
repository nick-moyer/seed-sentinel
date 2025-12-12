package store

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nick-moyer/seed-sentinel/models"
)

var db *sql.DB

func DB() *sql.DB {
	return db
}

func createTable(query string, tableName string) {
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create %s table: %v", tableName, err)
	}
}

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./data/sentinel.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	createTable(`
        CREATE TABLE IF NOT EXISTS sensors (
            id TEXT PRIMARY KEY,
            plant_name TEXT
        );`, "sensors") // id = device mac address

	createTable(`
        CREATE TABLE IF NOT EXISTS readings (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            sensor_id TEXT,
            moisture INTEGER,
            timestamp DATETIME,
            FOREIGN KEY(sensor_id) REFERENCES sensors(id)
        );`, "readings")

	log.Println("Database initialized (sentinel.db)")
}

func SaveSensor(id string, plantName string) {
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
