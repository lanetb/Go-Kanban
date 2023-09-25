package main

import (
	"html/template"
	"log"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Println("Application Started")
	// HandleFunc registers the handler function for the given pattern
	// in the DefaultServeMux.

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("../Client/html/index.html"))
		tmpl.Execute(w, nil)
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

	ConnectToDB()

	log.Fatal(http.ListenAndServe(":8000", nil))
}
