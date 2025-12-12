package store

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nick-moyer/seed-sentinel/models"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./data/sentinel.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Create the table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS readings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sensor_id TEXT,
		plant_name TEXT,
		moisture INTEGER,
		timestamp DATETIME
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	log.Println("Database initialized (seed.db)")
}

func SaveReading(payload models.SensorPayload) {
	stmt, err := db.Prepare("INSERT INTO readings(sensor_id, plant_name, moisture, timestamp) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Println("Database Error (Prepare):", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(payload.SensorID, payload.PlantName, payload.Moisture, time.Now())
	if err != nil {
		log.Println("Database Error (Insert):", err)
	} else {
		log.Println("Saved reading to DB")
	}
}
