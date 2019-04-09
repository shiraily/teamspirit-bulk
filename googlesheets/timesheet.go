package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/shiraily/teamspirit-bulk/model"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

var (
	sheetId       = os.Getenv("TS_SHEET_ID")
	sheetWorkTime = os.Getenv("TS_SHEET_WORK_TIME")
	sheetSetting  = os.Getenv("TS_SHEET_SETTING")
	clientSecret  = os.Getenv("TS_CLIENT_SECRET")
	strDayStart   = os.Getenv("TS_DAY_START")
)

func httpClient() (*http.Client, error) {
	conf, err := google.JWTConfigFromJSON([]byte(clientSecret), "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}

	return conf.Client(oauth2.NoContext), nil
}

func main() {
	timeDayStart, err := time.Parse("15:04", strDayStart)
	if err != nil {
		log.Fatal(err)
	}
	if noon, _ := time.Parse("15:04", "12:00"); !timeDayStart.Before(noon) {
		log.Fatal(err)
	}

	client, err := httpClient()
	if err != nil {
		log.Fatal(err)
	}
	sheetService, err := sheets.New(client)
	if err != nil {
		log.Fatalf("unable to retrieve Sheets Client %v", err)
	}
	ss := sheetService.Spreadsheets

	res, err := ss.Values.Get(sheetId, fmt.Sprintf("%s!A1", sheetSetting)).Do()
	if err != nil {
		log.Fatalf("unable to get Spreadsheets. %v", err)
	}
	if len(res.Values) == 0 {
		log.Fatalf("failed to get A1")
	}
	buf, ok := res.Values[0][0].(string)
	if !ok {
		log.Fatalf("failed to get first row: %s", res.Values[0][0])
	}
	firstRow, _ := strconv.Atoi(buf)
	res, err = ss.Values.Get(sheetId, fmt.Sprintf("%s!A%d:B%d", sheetWorkTime, firstRow, firstRow+3000)).Do()
	if err != nil {
		log.Fatal(err)
	}

	workTimes := [31]model.WorkTime{}
	curDay := 1
	curIsIn := true

	for _, row := range res.Values {
		isIn := row[0].(string) == "in"
		//FIXME
		if !(isIn == curIsIn) {
			continue
		}
		str := row[1].(string)
		log.Println(str)
		group := regexp.MustCompile(".* ([0-3][0-9]).*at ([0-9]{2}.[0-9]{2})(AM|PM)").FindStringSubmatch(str)
		log.Println(group)
		day, _ := strconv.Atoi(group[1])
		hhmm, _ := time.Parse("15:04", group[2])
		isAM := group[3] == "AM"
		var hhmmStr string
		//FIXME: some bugs
		if isAM && hhmm.Before(timeDayStart) {
			day -= 1
			hhmmStr = strconv.Itoa(hhmm.Hour() + 24)
		} else {
			hhmmStr = hhmm.Format("15:04")
		}
		curDay = day
		log.Printf("hhmm: %s", hhmmStr)
		if curIsIn {
			workTimes[curDay].StartTime = hhmm.String()
			curIsIn = !curIsIn
		} else {
			workTimes[curDay].EndTime = hhmm.String()
		}
	}

}

func getTime(a string) (string, error) {
	return "", nil
}
