package main

import (
	"log"

	"github.com/shiraily/teamspirit-bulk/teamspirit"
)

func main() {
	ts := teamspirit.NewTeamSpirit(teamspirit.DefaultDriver)
	if err := ts.Setup(); err != nil {
		log.Fatal(err)
	}

	// input 1
	if err := ts.InputWorkTime(1, "10:02", "19:07"); err != nil {
		log.Fatal(err)
	}
	log.Println("success to input time sheet")
}
