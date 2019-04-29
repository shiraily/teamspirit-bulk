package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/shiraily/teamspirit-bulk/config"
	"github.com/shiraily/teamspirit-bulk/googlesheets"
	"github.com/shiraily/teamspirit-bulk/teamspirit"
)

func main() {
	flag.String("", "", "config yaml file path")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.PrintDefaults()
		log.Fatal("config yaml file path is not given")
	}
	buf, err := ioutil.ReadFile(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = yaml.Unmarshal(buf, &config.Cfg)
	if err != nil {
		log.Fatalf("failed to unmarshal config yaml: %s", err)
	}

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
