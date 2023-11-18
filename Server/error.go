package main

import (
	"log"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Print("ERROR: ")
		log.Println(err)
		log.Fatal(msg)
	}
}
