package main

import (
	"html/template"
	"log"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

func GetProjects(w http.ResponseWriter ,r *http.Request){
	session, _ := store.Get(r, "session")
	CurrentUser := session.Values["CurrentUser"].(User) // type assert to User
	log.Println("Retrieving projects...")
	rows, err := db.Query("SELECT * FROM Project WHERE UserID=?", CurrentUser.ID)
	log.Println("here 4")
	if err != nil {
		log.Println(err)
	}
	log.Println("here 5")
	for rows.Next(){
		var project Project
		err := rows.Scan(&project.ID, &project.UserID, &project.Name, &project.LastUpdated)
		if err != nil {
			log.Println(err)
		}
		CurrentUser.Projects = append(CurrentUser.Projects, project) // access Projects field directly
	}
	session.Values["CurrentUser"] = CurrentUser // update session value
	session.Save(r, w)
	log.Println("Projects retrieved")
}

func CreateProject(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	projectName := r.FormValue("projectName")
	_, err := db.Exec("INSERT INTO Project (ProjectName, UserID) VALUES (?, ?)", projectName, CurrentUser.ID)
	if err != nil {
		log.Println(err)
	}
	GetProjects(w, r)
	tmpl, _ := template.ParseFiles("../Client/html/dashboard.html")
	err = tmpl.Execute(w, CurrentUser)
	if err != nil {
		log.Println(err)
	}
}

func OpenProjectHandler(w http.ResponseWriter, r *http.Request){
	log.Println("Opening project...")
	r.ParseForm()
	log.Println("Form: ", r.Form)
	val := r.FormValue("projectID")
	projectName := r.FormValue("projectName")
	//convert projectID to int
	projectID, err := strconv.Atoi(val)
	if err != nil {
		log.Println(err)
	}
	log.Println("Opening project: ", projectID)
	GetBoards(projectID)
	GetTasks(projectID)
	for i, board := range Boards{
		for _, task := range Tasks{
			if board.ID == task.BoardID{
				Boards[i].Tasks = append(Boards[i].Tasks, task)
			}
		}
	}
	tmpl, _ := template.ParseFiles("../Client/html/project.html")
	data := struct{
		ProjectName string
		User User
		Boards []Board
	}{
		ProjectName: projectName,
		User: CurrentUser,
		Boards: Boards,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
	log.Println("Project opened")
}

func GetBoards(projectID int){
	log.Println("Retrieving boards...")
	rows, err := db.Query("SELECT * FROM Board WHERE ProjectID=?", projectID)
	if err != nil {
		log.Println(err)
	}
	for rows.Next(){
		var board Board
		err := rows.Scan(&board.ID, &board.ProjectID, &board.UserID, &board.Name)
		if err != nil {
			log.Println(err)
		}
		Boards = append(Boards, board)
	}
	log.Println("Boards retrieved")
}

func GetTasks(projectID int){
	log.Println("Retrieving tasks...")
	rows, err := db.Query("SELECT * FROM Task WHERE ProjectID=?", projectID)
	if err != nil {
		log.Println(err)
	}
	for rows.Next(){
		var task Task
		err := rows.Scan(&task.ID, &task.BoardID, &task.ProjectID, &task.UserID, &task.Name, &task.Description, &task.Type)
		if err != nil {
			log.Println(err)
		}
		Tasks = append(Tasks, task)
	}
	log.Println("Tasks retrieved")
}