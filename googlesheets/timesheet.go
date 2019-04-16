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

	dayWorkTime := map[int]model.WorkTime{}
	var lastDay int
	var lastHhmmStr string

	for _, row := range res.Values {
		isIn := row[0].(string) == "in"
		str := row[1].(string)
		group := regexp.MustCompile(
			".* ([0-3][0-9]), [0-9]{4} at ([0-9]{2}.[0-9]{2})(AM|PM)",
		).FindStringSubmatch(str)

		day, _ := strconv.Atoi(group[1])
		hhmmStr := group[2]
		hhmm, _ := time.Parse("15:04", hhmmStr)
		hour := hhmm.Hour()
		isAM := group[3] == "AM"
		if isAM && hhmm.Before(timeDayStart) {
			day -= 1
			hhmmStr = fmt.Sprintf("%02d:%02d", hour+24, hhmm.Minute())
		} else if !isAM && hour != 12 {
			hhmmStr = fmt.Sprintf("%02d:%02d", hour+12, hhmm.Minute())
			hhmm, _ = time.Parse("15:04", hhmmStr)
		}

		if hhmm.After(timeDayStart) {
			if isIn && dayWorkTime[day].StartTime == "" {
				dayWorkTime[lastDay] = model.WorkTime{
					Day:       lastDay,
					StartTime: dayWorkTime[lastDay].StartTime,
					EndTime:   lastHhmmStr,
				}
				dayWorkTime[day] = model.WorkTime{
					Day:       day,
					StartTime: hhmmStr,
				}
			} else if dayWorkTime[day].StartTime == "" {
				continue
			}
		}

		if !isIn {
			lastDay = day
			lastHhmmStr = hhmmStr
		}
	}

	dayWorkTime[lastDay] = model.WorkTime{
		Day:       lastDay,
		StartTime: dayWorkTime[lastDay].StartTime,
		EndTime:   lastHhmmStr,
	}

	for day := 1; day <= 31; day++ {
		if wt, ok := dayWorkTime[day]; ok {
			log.Printf("%02d: %s - %s", day, wt.StartTime, wt.EndTime)
		}
	}
	// TODO stash 0 day
}
