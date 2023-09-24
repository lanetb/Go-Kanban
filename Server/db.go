package main

import (
	"database/sql"
	"log"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

// DB is a global variable that holds the connection to the database
func ConnectToDB(){
	log.Println("Connecting to database...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DB, err := sql.Open("mysql", os.Getenv("DSN"))

    if err != nil {

        log.Fatalf("failed to connect: %v", err)

    }

    defer DB.Close()

    if err := DB.Ping(); err != nil {

        log.Fatalf("failed to ping: %v", err)

    }

    log.Println("Successfully connected to PlanetScale!")
}

func RegistraitionAuthHandler(w http.ResponseWriter, r *http.Request) {
	// ParseForm parses the raw query from the URL and updates r.Form.
	// For POST requests, it also parses the request body as a form and puts the results into both r.PostForm and r.Form.
	r.ParseForm()
	// FormValue returns the first value for the named component of the query.
	// POST and PUT body parameters take precedence over URL query string values.
	// FormValue calls ParseMultipartForm and ParseForm if necessary and ignores any errors returned by these functions.
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	// Check if username is already taken
	var user string
	// check if password, email and username are valid
	// username must be 4 or more characters and consitst of only letters and numbers
	// password must be 8 or more characters and consist of letters, numbers and special characters (e.g. !@#$%^&*)'
	// password must contain at least one uppercase letter, one lowercase letter, one number and one special character
	// email must be valid email address
	// if not, redirect to registration page


	err := DB.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)
	switch {
		case err == sql.ErrNoRows:	
			// Username is not taken
			// Insert user into database
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				log.Panic(err)
			}
			password = string(hashedPassword)
			
			_, err = DB.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?)", username, password, email)
			if err != nil {
				log.Println(err)
			}
			// Redirect to login page
			log.Println("User created")
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		case err != nil:
			log.Println(err)
		default:
			// Username is taken
			// Redirect to registration page
			log.Println("Username taken")
			http.Redirect(w, r, "/register", http.StatusMovedPermanently)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request){
	// ParseForm parses the raw query from the URL and updates r.Form.
	// For POST requests, it also parses the request body as a form and puts the results into both r.PostForm and r.Form.
	r.ParseForm()
	// FormValue returns the first value for the named component of the query.
	// POST and PUT body parameters take precedence over URL query string values.
	// FormValue calls ParseMultipartForm and ParseForm if necessary and ignores any errors returned by these functions.
	username := r.FormValue("username")
	password := r.FormValue("password")
	// Check if username is already taken
	var user string
	var pass string
	err := DB.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&user, &pass)
	switch {
		case err == sql.ErrNoRows:	
			// Username is not taken
			// Redirect to login page
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		case err != nil:
			log.Println(err)
		default:
			// Username is taken
			// Check password
			err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
			if err != nil {
				log.Println(err)
				// Redirect to login page
				http.Redirect(w, r, "/", http.StatusMovedPermanently)
			} else {
				// Redirect to home page
				http.Redirect(w, r, "/home", http.StatusMovedPermanently)
			}
	}
}