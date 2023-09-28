package main

import (
	"html/template"
	"log"
	"net/http"
)

func ReturnToDashHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		log.Println(err)
	}
	if session.Values["CurrentUser"] == nil {
		tmpl := template.Must(template.ParseFiles("../Client/html/index.html"))
		tmpl.Execute(w, nil)
	} else {
		tmpl := template.Must(template.ParseFiles("../Client/html/indexsigned.html"))
		data := struct {
			CurrUser User
			Projects map[int]Project
		}{
			CurrUser: session.Values["CurrentUser"].(User),
			Projects: session.Values["CurrentUser"].(User).Projects,
		}
		tmpl.Execute(w, data)
	}
}

func SignoutHandler(w http.ResponseWriter, r *http.Request){
	session, err := store.Get(r, "session")
	if err != nil {
		log.Println(err)
	}
	session.Values["CurrentUser"] = nil
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		log.Println(err)
	}
	tmpl := template.Must(template.ParseFiles("../Client/html/index.html"))
	tmpl.Execute(w, nil)
}