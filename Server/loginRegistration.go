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

	if len(password) >= 8 && len(password) <= 100 {
		passLength = true
	}

	err := db.QueryRow("SELECT username FROM User WHERE username=?", username).Scan(&user)
	switch {
		case err == sql.ErrNoRows:
			if usernameValid && usernameLength{
			    if passLower && passUpper  && passNumber && passSpecial && passLength && passNoSpace {
					hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
					handleError(err, "Error hashing password")
					password = string(hashedPassword)
					_, err = db.Exec("INSERT INTO User (username, password, email) VALUES (?, ?, ?)", username, password, email)
					handleError(err, "Error inserting user")
					log.Println("User created")
					Login(username, password, w, r)
				} else {
					tmpl := template.Must(template.ParseFiles("../Client/html/registerbadpassword.html"))
					tmpl.Execute(w, nil)
				}
			} else {
				tmpl := template.Must(template.ParseFiles("../Client/html/registerbaduser.html"))
				tmpl.Execute(w, nil)
			}
		case err != nil:
			log.Println(err)
		default:
			tmpl := template.Must(template.ParseFiles("../Client/html/registerusertaken.html"))
			tmpl.Execute(w, nil)
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
		template, _ := template.ParseFiles("../Client/html/loginerror.html")
		template.Execute(w, nil)
	} else{
		switch {
			case err == sql.ErrNoRows:
				template, _ := template.ParseFiles("../Client/html/loginerror.html")
				template.Execute(w, nil)
			case err != nil:
				template, _ := template.ParseFiles("../Client/html/loginerroruserpass.html")
				template.Execute(w, nil)
			default:
				err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
				if err != nil {
					template, _ := template.ParseFiles("../Client/html/loginerroruserpass.html")
					template.Execute(w, nil)
				} else {
					log.Print("Logged in")
					var currentUser User
					currentUser = User{
						Username: user,
						ID: ID,
						Projects: make(map[int]Project),
					}
					session, _ := store.Get(r, "session")
					session.Options.MaxAge = 60 * 60 * 24
					session.Values["CurrentUser"] = currentUser
					log.Println("User: ", currentUser.Username)
					session.Save(r, w)

					currentUser = GetProjects(w, r, session)
					session.Save(r, w)
					tmpl, _ := template.ParseFiles("../Client/html/dashboard.html")
					data := struct{
						User User 
						Projects map[int]Project
					}{
						User: currentUser,
						Projects: currentUser.Projects,
					}
					err = tmpl.Execute(w, data)
					handleError(err, "Error executing template")
				}
		}
	}
}