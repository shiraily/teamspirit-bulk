package main

import (
	"log"

	"github.com/shiraily/teamspirit-bulk/googlesheets"
	"github.com/shiraily/teamspirit-bulk/teamspirit"
)

func main() {
	httpClient, err := googlesheets.NewGoogleSheetsClient()
	if err != nil {
		log.Fatal(err)
	}
	sheet, err := googlesheets.NewTimeSheet(httpClient)
	if err != nil {
		log.Fatal(err)
	}
	if err := sheet.Setup(); err != nil {
		log.Fatal(err)
	}
	workTimes, err := sheet.GetWorkTimes()
	if err != nil {
		log.Fatal(err)
	}

	ts := teamspirit.NewTeamSpirit(teamspirit.DefaultDriver)
	if err := ts.Setup(); err != nil {
		log.Fatal(err)
	}
	err = ts.BulkInput(workTimes)
	ts.Driver.Stop() // TODO consider when should stop
	if err != nil {
		log.Fatal(err)
	}
	log.Println("success to input time sheet")
}
