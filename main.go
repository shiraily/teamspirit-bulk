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
	if err := ts.InputWorkTime(t.page, 1, "10:02", false); err != nil {
		log.Fatal(err)
	}
	log.Println("success to input time sheet")
}
