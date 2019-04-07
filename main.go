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
	wts := []teamspirit.WorkTime{
		{
			Day:       1,
			StartTime: "10:02",
			EndTime:   "19:07",
		},
		{
			Day:       2,
			StartTime: "10:03",
			EndTime:   "19:09",
		},
	}
	err := ts.BulkInput(wts)
	ts.Driver.Stop()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("success to input time sheet")
}
