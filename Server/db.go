package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"unicode"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)
var db *sql.DB
type User struct{
	Username string
	ID int
}
var CurrentUser User
type Project struct{
	ID int
	UserID int
	Name string
	LastUpdated string
}

type Board struct{
	ID int
	ProjectID int
	UserID int
	Name string
}

type Task struct{
	ID int
	BoardID int
	ProjectID int
	UserID int
	Name string
	Description string
	Type string
}

var Projects []Project
var Boards []Board
var Tasks []Task
//    db is a global variable that holds the connection to the database
func ConnectToDB(){
	log.Println("Connecting to database...")
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

    db, err = sql.Open("mysql", os.Getenv("DSN"))

    if err != nil {
        log.Fatalf("failed to connect: %v", err)
    }

    if err := db.Ping(); err != nil {
        log.Fatalf("failed to ping: %v", err)
    }

    log.Println("Successfully connected to PlanetScale!")
}

func RegistraitionAuthHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(db)
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

	if len(password) >= 8 && len(password) <= 50 {
		passLength = true
	}

	err := db.QueryRow("SELECT username FROM User WHERE username=?", username).Scan(&user)
	switch {
		case err == sql.ErrNoRows:
			if usernameValid && usernameLength && passLower && passUpper  && passNumber && passSpecial && passLength && passNoSpace {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				if err != nil {
					log.Println(err)
				}
				password = string(hashedPassword)
				log.Println("hash: ", password)
				_, err = db.Exec("INSERT INTO User (username, password, email) VALUES (?, ?, ?)", username, password, email)
				if err != nil {
					log.Println(err)
				}
				log.Println("User created")
				//Login(username, password, w, r)
			} else {
				log.Println("Password invalid")
			}
		case err != nil:
			log.Println(err)
		default:
			log.Println("Username taken")
			Login(username, password, w, r)
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
	log.Print("Logging in...")
	err := db.QueryRow("SELECT username, password FROM User WHERE username=?", username).Scan(&user, &pass)
	if err != nil {
		log.Println(err)
	}
	switch {
		case err == sql.ErrNoRows:
			log.Println(err)
		case err != nil:
			log.Println(err)
		default:
			err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
			if err != nil {
				log.Println(err)
			} else {
				log.Print("Logged in")
				CurrentUser = User{
					Username: user,
					ID: 1,
				}
				GetProjects()
				tmpl, _ := template.ParseFiles("../Client/html/dashboard.html")
				log.Println("Projects: ", Projects)
				data := struct{
					User User
					Projects []Project
				}{
					User: CurrentUser,
					Projects: Projects,
				}
				err = tmpl.Execute(w, data)
				if err != nil {
					log.Println(err)
				}
			}
	}
}

func GetProjects(){
	log.Println("Retrieving projects...")
	rows, err := db.Query("SELECT * FROM Project WHERE UserID=?", CurrentUser.ID)
	if err != nil {
		log.Println(err)
	}
	for rows.Next(){
		var project Project
		err := rows.Scan(&project.ID, &project.UserID, &project.Name, &project.LastUpdated)
		if err != nil {
			log.Println(err)
		}
		Projects = append(Projects, project)
	}
	log.Println("Projects retrieved")
}

func CreateProject(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	projectName := r.FormValue("projectName")
	_, err := db.Exec("INSERT INTO Project (ProjectName, UserID) VALUES (?, ?)", projectName, CurrentUser.ID)
	if err != nil {
		log.Println(err)
	}
	GetProjects()
	tmpl, _ := template.ParseFiles("../Client/html/dashboard.html")
	err = tmpl.Execute(w, CurrentUser)
	if err != nil {
		log.Println(err)
	}
}

func OpenProjectHandler(w http.ResponseWriter, r *http.Request){
	log.Println("Opening project...")
	r.ParseForm()
	val := r.FormValue("projectID")
	projectName := r.FormValue("Name")
	//convert projectID to int
	projectID, err := strconv.Atoi(val)
	if err != nil {
		log.Println(err)
	}
	log.Println("Opening project: ", projectID)
	GetBoards(projectID)
	GetTasks(projectID)
	tmpl, _ := template.ParseFiles("../Client/html/project.html")
	data := struct{
		ProjectName string
		User User
		Boards []Board
		Tasks []Task
	}{
		ProjectName: projectName,
		User: CurrentUser,
		Boards: Boards,
		Tasks: Tasks,
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