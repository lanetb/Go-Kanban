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
		tmpl := template.Must(template.ParseFiles("../Client/html/login.html"))
		tmpl.Execute(w, nil)
	}

	http.HandleFunc("/", h1)
	//var db = ConnectToDB()

	log.Fatal(http.ListenAndServe(":8000", nil))
}
