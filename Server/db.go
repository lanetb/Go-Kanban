package main

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)
var db *sql.DB
var CurrentUser User
var Projects []Project
var Boards []Board
var Tasks []Task
//    db is a global variable that holds the connection to the database
func ConnectToDB(){
	log.Println("Connecting to database...")
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

    db, err = sql.Open("mysql", os.Getenv("DSN"))

    if err != nil {
        log.Fatalf("failed to connect: %v", err)
    }

    if err := db.Ping(); err != nil {
        log.Fatalf("failed to ping: %v", err)
    }
    log.Println("Successfully connected to PlanetScale!")
}