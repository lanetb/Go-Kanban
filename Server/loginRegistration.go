package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

func RegistraitionAuthHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(db)
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	var user string
	var usernameValid bool = true
	var usernameLength bool = true
	var passLower, passUpper, passNumber, passSpecial, passLength, passNoSpace bool = false, false, false, false, false, true

	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			usernameValid = false
		}
	}

	if len(username) <= 3 || len(username) >= 30 {
		usernameLength = false
	}

	for _, char := range password {
		switch{
			case unicode.IsLower(char):
				passLower = true
			case unicode.IsUpper(char):
				passUpper = true
			case unicode.IsNumber(char):
				passNumber = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				passSpecial = true
			case unicode.IsSpace(char):
				passNoSpace = false
		}
	}

	if len(password) >= 8 && len(password) <= 50 {
		passLength = true
	}

	err := db.QueryRow("SELECT username FROM User WHERE username=?", username).Scan(&user)
	switch {
		case err == sql.ErrNoRows:
			if usernameValid && usernameLength && passLower && passUpper  && passNumber && passSpecial && passLength && passNoSpace {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				if err != nil {
					log.Println(err)
				}
				password = string(hashedPassword)
				log.Println("hash: ", password)
				_, err = db.Exec("INSERT INTO User (username, password, email) VALUES (?, ?, ?)", username, password, email)
				if err != nil {
					log.Println(err)
				}
				log.Println("User created")
				//Login(username, password, w, r)
			} else {
				log.Println("Password invalid")
			}
		case err != nil:
			log.Println(err)
		default:
			log.Println("Username taken")
			Login(username, password, w, r)
	}

}

func LoginAuthHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	Login(username, password, w, r)
}

func Login(username string, password string, w http.ResponseWriter, r *http.Request){
	var user string
	var pass string
	var ID int
	log.Print("Logging in...")
	err := db.QueryRow("SELECT UserID, Username, Password FROM User WHERE Username=?", username).Scan(&ID, &user, &pass)
	if err != nil {
		log.Println(err)
	}
	switch {
		case err == sql.ErrNoRows:
			log.Println(err)
		case err != nil:
			log.Println(err)
		default:
			err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
			if err != nil {
				log.Println(err)
			} else {
				log.Print("Logged in")
				session, _ := store.Get(r, "session")
				log.Print("here 1")
				session.Values["CurrentUser"] = CurrentUser
				log.Print("here 2")
				currentUser := session.Values["CurrentUser"].(User)
				currentUser.Username = user
				currentUser.ID = ID
				session.Values["CurrentUser"] = currentUser
				log.Print("here 3")
				session.Save(r, w)
				GetProjects(w, r)
				CurrentUser := session.Values["CurrentUser"]
				tmpl, _ := template.ParseFiles("../Client/html/dashboard.html")
				log.Println("Projects: ", Projects)
				data := struct{
					User User 
					Projects []Project
				}{
					User: CurrentUser.(User),
					Projects: CurrentUser.(User).Projects,
				}
				err = tmpl.Execute(w, data)
				if err != nil {
					log.Println(err)
				}
			}
	}
}