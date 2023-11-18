package main

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	
)
var db *sql.DB
//    db is a global variable that holds the connection to the database
func ConnectToDB(){
	log.Println("Connecting to database...")
	err := godotenv.Load()

	handleError(err, "Error loading .env file")

    db, err = sql.Open("mysql", os.Getenv("DSN"))

    handleError(err, "Error opening database")

    if err := db.Ping(); err != nil {
        log.Fatalf("failed to ping: %v", err)
    }
    log.Println("Successfully connected to PlanetScale!")
}