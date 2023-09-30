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
		if len(r.Cookies()) == 0 {
			log.Println("No cookies")
			tmpl := template.Must(template.ParseFiles("../Client/html/index.html"))
			tmpl.Execute(w, nil)
		} else {
			log.Println("Cookies")
			session, err := store.Get(r, "session")
			if err != nil {
				log.Println(err)
			}
			log.Println(session.Values["CurrentUser"].(User).Username)
			log.Println(session.Values["CurrentUser"].(User).Projects)
				data := struct{
					CurrUser User 
					Projects map[int]Project
				}{
					CurrUser: session.Values["CurrentUser"].(User),
					Projects: session.Values["CurrentUser"].(User).Projects,
				}
			tmpl := template.Must(template.ParseFiles("../Client/html/indexsigned.html"))
			tmpl.Execute(w, data)
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
	http.HandleFunc("/createProject/", CreateProjectHandler)
	http.HandleFunc("/createBoard/", CreateBoardHandler)
	http.HandleFunc("/createTask/", CreateTaskHandler)
	http.HandleFunc("/returnToDash/", ReturnToDashHandler)
	http.HandleFunc("/signout/", SignoutHandler)
	http.HandleFunc("/deleteProject/", DeleteProjectHandler)
	http.HandleFunc("/deleteBoard/", DeleteBoardHandler)
	http.HandleFunc("/onDragEnd/", HandleDragEnd)

	ConnectToDB()

	log.Fatal(http.ListenAndServe(":8000", context.ClearHandler(http.DefaultServeMux)))
}
