package main

import (
	"database/sql"
	"log"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

// DB is a global variable that holds the connection to the database
func ConnectToDB() *sql.DB{
	log.Println("Connecting to database...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("mysql", os.Getenv("DSN"))

    if err != nil {

        log.Fatalf("failed to connect: %v", err)

    }

    defer db.Close()

    if err := db.Ping(); err != nil {

        log.Fatalf("failed to ping: %v", err)

    }

    log.Println("Successfully connected to PlanetScale!")
	return db
}
