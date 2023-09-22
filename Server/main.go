package main

import (
	"html/template"
	"log"
	"os"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Application Started")
	// HandleFunc registers the handler function for the given pattern
	// in the DefaultServeMux.

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("../Client/HTML/index.html"))
		tmpl.Execute(w, nil)
	}

	http.HandleFunc("/", h1)
	
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

	log.Fatal(http.ListenAndServe(":8000", nil))
}
