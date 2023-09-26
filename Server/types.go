package main

type User struct{
	Username string
	ID 		 int
	Projects []Project
}

type Project struct{
	ID 			int
	UserID 		int
	Name 		string
	LastUpdated string
	Boards      []Board
}

type Board struct{
	ID 			int
	ProjectID 	int
	UserID 		int
	Name 		string
	Tasks       []Task
}

type Task struct{
	ID 			int
	BoardID 	int
	ProjectID 	int
	UserID 		int
	Name 		string
	Description string
	Type 		string
}