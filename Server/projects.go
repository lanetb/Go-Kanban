package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

func GetProjects(w http.ResponseWriter ,r *http.Request, session *sessions.Session) User{
	// type assert to User
	CurrentUser := session.Values["CurrentUser"].(User)
	log.Println("Retrieving projects...")
	rows, err := db.Query("SELECT * FROM Project WHERE UserID=?", CurrentUser.ID)
	if err != nil {
		log.Println(err)
	}
	projects := make(map[int]Project)
	for rows.Next(){
		var project Project
		err := rows.Scan(&project.ID, &project.UserID, &project.Name, &project.LastUpdated)
		if err != nil {
			log.Println(err)
		}
		projects[project.ID] = project // access Projects field directly
	}
	CurrentUser.Projects = projects // fix type assertion error
	session.Values["CurrentUser"] = CurrentUser // update session value
	log.Println("Projects retrieved")
	return CurrentUser
}	

func CreateProjectHandler(w http.ResponseWriter, r *http.Request){
	log.Println("Creating project...")
	r.ParseForm()
	session, _ := store.Get(r, "session")
	CurrentUser := session.Values["CurrentUser"].(User)
	projectName := r.FormValue("projectName")
	// insert project into database
	stmt, err := db.Prepare("INSERT INTO Project (UserID, ProjectName) VALUES (?, ?)")
	if err != nil {
		log.Println(err)
	}
	res, err := stmt.Exec(CurrentUser.ID, projectName)
	if err != nil {
		log.Println(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	BuildBoard(session, w, r, int(id), CurrentUser, "Backlog")
	BuildBoard(session, w, r, int(id), CurrentUser, "In Progress")
	BuildBoard(session, w, r, int(id), CurrentUser, "Finished")
	log.Println("Project created")
	// update session value
	CurrentUser.Projects[int(id)] = Project{ID: int(id), UserID: CurrentUser.ID, Name: projectName}
	session.Values["CurrentUser"] = CurrentUser
	session.Save(r, w)
	log.Println("Project created")

	tmpl, _ := template.ParseFiles("../Client/html/dashboard.html")
	data := struct{
		User User 
		Projects map[int]Project
	}{
		User: session.Values["CurrentUser"].(User),
		Projects: session.Values["CurrentUser"].(User).Projects,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

func CreateBoardHandler(w http.ResponseWriter, r *http.Request){
	log.Println("Creating board...")
	r.ParseForm()
	session, _ := store.Get(r, "session")
	CurrentUser := session.Values["CurrentUser"].(User)
	projectID, err := strconv.Atoi(r.FormValue("projectID"))
	boardName := r.FormValue("boardName")
	log.Println("ProjectID: ", projectID)
	if err != nil {
		log.Println(err)
	}
	// insert board into database
	BuildBoard(session, w, r, projectID, CurrentUser, boardName)
	tmpl, _ := template.ParseFiles("../Client/html/project.html")
	data := struct{
		ProjectName string
		User User
		Boards []Board
	}{
		ProjectName: CurrentUser.Projects[projectID].Name,
		User: session.Values["CurrentUser"].(User),
		Boards: session.Values["CurrentUser"].(User).Projects[projectID].Boards,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

func BuildBoard(session *sessions.Session, w http.ResponseWriter, r *http.Request, projectID int, CurrentUser User, boardName string) {
	stmt, err := db.Prepare("INSERT INTO Board (ProjectID, UserID, BoardName) VALUES (?, ?, ?)")
	if err != nil {
		log.Println(err)
	}
	res, err := stmt.Exec(projectID, CurrentUser.ID, boardName)
	if err != nil {
		log.Println(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	log.Println("Board created")
	// update session value
	project := CurrentUser.Projects[projectID]
	project.Boards = append(project.Boards, Board{ID: int(id), ProjectID: projectID, UserID: CurrentUser.ID, Name: boardName})
	CurrentUser.Projects[projectID] = project
	session.Values["CurrentUser"] = CurrentUser
	session.Save(r, w)
	log.Println("Board created")
}


func OpenProjectHandler(w http.ResponseWriter, r *http.Request){
	log.Println("Opening project...")
	r.ParseForm()
	log.Println("Form: ", r.Form)
	val := r.FormValue("projectID")
	projectName := r.FormValue("projectName")
	//convert projectID to int
	projectID, err := strconv.Atoi(val)
	session, _ := store.Get(r, "session")
	CurrentUser := session.Values["CurrentUser"].(User)
	if err != nil {
		log.Println(err)
	}
	log.Println("Opening project: ", projectID)
	GetBoards(w, r, projectID)
	GetTasks(w, r, projectID)
	log.Println("Project opened: ", projectID)
	log.Println("Project name: ", projectName)

	tmpl, _ := template.ParseFiles("../Client/html/project.html")
	session.Values["CurrentUser"] = CurrentUser
	data := struct{
		ProjectName string
		User User
		Project Project
		Boards []Board
	}{
		ProjectName: projectName,
		User: session.Values["CurrentUser"].(User),
		Project: session.Values["CurrentUser"].(User).Projects[projectID],
		Boards: session.Values["CurrentUser"].(User).Projects[projectID].Boards,
	}
	session.Save(r, w)
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
	log.Println("Project opened")
}

func GetBoards(w http.ResponseWriter ,r *http.Request, projectID int) {
	log.Println("Retrieving boards...")
	rows, err := db.Query("SELECT * FROM Board WHERE ProjectID=?", projectID)
	session, _ := store.Get(r, "session")
	CurrentUser := session.Values["CurrentUser"].(User)
	CurrProject := CurrentUser.Projects[projectID]
	var tempBoard []Board
	if err != nil {
		log.Println(err)
	}
	for rows.Next(){
		var board Board
		err := rows.Scan(&board.ID, &board.ProjectID, &board.UserID, &board.Name)
		if err != nil {
			log.Println(err)
		}
		tempBoard = append(tempBoard, board)
	}
	log.Println(tempBoard)
	CurrProject.Boards = tempBoard
	CurrentUser.Projects[projectID] = CurrProject
	session.Values["CurrentUser"] = CurrentUser
	session.Save(r, w)
	log.Println("Boards retrieved")
}

func GetTasks(w http.ResponseWriter ,r *http.Request, projectID int) {
	rows, err := db.Query("SELECT * FROM Task WHERE ProjectID=?", projectID)
	session, _ := store.Get(r, "session")
	CurrentUser := session.Values["CurrentUser"].(User)
	CurrProject := CurrentUser.Projects[projectID]
	if err != nil {
		log.Println(err)
	}
	for rows.Next(){
		var task Task
		err := rows.Scan(&task.ID, &task.BoardID, &task.ProjectID, &task.UserID, &task.Name, &task.Description, &task.Type)
		if err != nil {
			log.Println(err)
		}
		for i, board := range CurrProject.Boards{
			if board.ID == task.BoardID{
				CurrProject.Boards[i].Tasks = append(CurrProject.Boards[i].Tasks, task)
			}
		}
	}
	CurrentUser.Projects[projectID] = CurrProject
	session.Values["CurrentUser"] = CurrentUser
	log.Println("Tasks retrieved")
}