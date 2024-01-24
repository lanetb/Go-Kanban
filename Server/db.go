package main

import (
	"database/sql"
	"log"
	"os"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	
)
var db *sql.DB
//    db is a global variable that holds the connection to the database
func ConnectToDB(){
	log.Println("Connecting to database...")
	err := godotenv.Load()

	handleError(err, "Error loading .env file")

	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASS"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_HOST"),
		DBName:               os.Getenv("DB_NAME"),
	}

    db, err = sql.Open("mysql", cfg.FormatDSN())

    handleError(err, "Error opening database")

    if err := db.Ping(); err != nil {
        log.Fatalf("failed to ping: %v", err)
    }
    log.Println("Successfully connected to PlanetScale!")
	log.Println("Go to http://localhost:8080/ to view the app")
}