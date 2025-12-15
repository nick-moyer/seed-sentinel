package store

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
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
