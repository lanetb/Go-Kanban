package main

import (
	"html/template"
	"net/http"
)

func ReturnToDashHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	handleError(err, "Error getting session")
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
	handleError(err, "Error getting session")
	session.Values["CurrentUser"] = nil
	session.Options.MaxAge = -1
	handleError(session.Save(r, w), "Error saving session")
	tmpl := template.Must(template.ParseFiles("../Client/html/index.html"))
	tmpl.Execute(w, nil)
}