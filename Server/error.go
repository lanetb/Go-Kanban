package main;

import (
	"log"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatal(msg)
	}
}
