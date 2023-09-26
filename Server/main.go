package main

import (
	"html/template"
	"log"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"encoding/gob"
)

var store = sessions.NewCookieStore([]byte("secret-key"))

func main() {
	log.Println("Application Started")
	// HandleFunc registers the handler function for the given pattern
	// in the DefaultServeMux.
	gob.Register(User{})
	gob.Register(Project{})
	gob.Register(Board{})
	gob.Register(Task{})

	h1 := func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		var User User
		if session.Values["CurrentUser"] == nil {
			log.Println("herio 1")
			session.Values["CurrentUser"] = User
			err := session.Save(r, w)
			if err != nil {
				log.Println(err)
			}
			tmpl := template.Must(template.ParseFiles("../Client/html/index.html"))
			tmpl.Execute(w, nil)
		} else {
			log.Println("herio 2")
			tmpl := template.Must(template.ParseFiles("../Client/html/index.html"))
			tmpl.Execute(w, nil)
		}
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("../Client/html/register.html"))
		tmpl.Execute(w, nil)
	}
	h3 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("../Client/html/login.html"))
		tmpl.Execute(w, nil)
	}

	fs := http.FileServer(http.Dir("../Client/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", h1)
	http.HandleFunc("/swap-reg/", h2)
	http.HandleFunc("/swap-log/", h3)
	http.HandleFunc("/register/", RegistraitionAuthHandler)
	http.HandleFunc("/login/", LoginAuthHandler)
	http.HandleFunc("/openProject/", OpenProjectHandler)

	ConnectToDB()

	log.Fatal(http.ListenAndServe(":8000", context.ClearHandler(http.DefaultServeMux)))
}
