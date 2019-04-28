package googlesheets

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

type TimeSheet struct {
	sheetService *sheets.SpreadsheetsService
	sheetID      string
	timeDayStart time.Time
	firstRow     int
}

func NewTimeSheet(client *http.Client) (*TimeSheet, error) {
	sheetService, err := sheets.New(client)
	if err != nil {
		log.Fatalf("unable to retrieve Sheets Client %v", err)
	}
	ss := sheetService.Spreadsheets

	return &TimeSheet{
		sheetService: ss,
	}, nil
}

func (t *TimeSheet) Setup() error {
	t.sheetID = sheetId
	timeDayStart, err := time.Parse("15:04", strDayStart)
	if err != nil {
		return fmt.Errorf("irregal time string format: %s", err)
	}
	if noon, _ := time.Parse("15:04", "12:00"); !timeDayStart.Before(noon) {
		return fmt.Errorf("start work time should be before noon")
	}
	t.timeDayStart = timeDayStart

	// first row
	res, err := t.sheetService.Values.Get(sheetId, fmt.Sprintf("%s!A1", sheetSetting)).Do()
	if err != nil {
		return fmt.Errorf("unable to get settings sheets. %s", err)
	}
	if len(res.Values) == 0 {
		return fmt.Errorf("failed to get A1 cell. sheet=%s", sheetSetting)
	}
	buf, ok := res.Values[0][0].(string)
	if !ok {
		return fmt.Errorf("invalid first row. value=%s", res.Values[0][0])
	}
	firstRow, err := strconv.Atoi(buf)
	if err != nil {
		return fmt.Errorf("firstRow is not integer: %s", err)
	}
	if firstRow < 1 {
		return fmt.Errorf("firstRow is less than 1. actual=%d", firstRow)
	}
	t.firstRow = firstRow

	return nil
}

func (t *TimeSheet) GetWorkTimes() ([]model.WorkTime, error) {
	res, err := t.sheetService.Values.Get(
		sheetId,
		fmt.Sprintf("%s!A%d:B%d", sheetWorkTime, t.firstRow, t.firstRow+3000),
	).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get work time data: %s", err)
	}

	dayWorkTime := map[int]model.WorkTime{}
	var lastDay int
	var lastHhmmStr string
	for _, row := range res.Values {
		isIn := row[0].(string) == "in"
		day, hhmmStr, hhmm := t.splitTimeFormat(row[1].(string))

		if hhmm.After(t.timeDayStart) {
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

	var workTimes []model.WorkTime
	for day := 1; day <= 31; day++ {
		if wt, ok := dayWorkTime[day]; ok {
			workTimes = append(workTimes, dayWorkTime[day])
			log.Printf("%02d: %s - %s", day, wt.StartTime, wt.EndTime)
		}
	}
	return workTimes, nil
}

func (t *TimeSheet) splitTimeFormat(buf string) (int, string, time.Time) {
	group := regexp.MustCompile(
		".* ([0-3][0-9]), [0-9]{4} at ([0-9]{2}.[0-9]{2})(AM|PM)",
	).FindStringSubmatch(buf)
	day, _ := strconv.Atoi(group[1])
	hhmmStr := group[2]
	isAM := group[3] == "AM"
	hhmm, _ := time.Parse("15:04", hhmmStr)
	hour := hhmm.Hour()
	if isAM && hhmm.Before(t.timeDayStart) {
		day -= 1
		hhmmStr = fmt.Sprintf("%02d:%02d", hour+24, hhmm.Minute())
	} else if !isAM && hour != 12 {
		hhmmStr = fmt.Sprintf("%02d:%02d", hour+12, hhmm.Minute())
		hhmm, _ = time.Parse("15:04", hhmmStr)
	}
	return day, hhmmStr, hhmm
}

func NewGoogleSheetsClient() (*http.Client, error) {
	conf, err := google.JWTConfigFromJSON([]byte(clientSecret), "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, fmt.Errorf("failed to read jwt: %s", err)
	}
	return conf.Client(oauth2.NoContext), nil
}
