package main

import (
	"fmt"
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
		username := r.FormValue("username")
		fmt.Println(username)
	}
	h3 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("../Client/html/register.html"))
		tmpl.Execute(w, nil)
	}
	h4 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("../Client/html/login.html"))
		tmpl.Execute(w, nil)
	}
	http.HandleFunc("/", h1)
	http.HandleFunc("/register/", h2)
	http.HandleFunc("/swap-reg/", h3)
	http.HandleFunc("/login/", h4)
	ConnectToDB()

	log.Fatal(http.ListenAndServe(":8000", nil))
}
