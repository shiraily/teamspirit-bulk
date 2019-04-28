package googlesheets

import (
	"net/http"
	"testing"
	"time"
)

func TestSplitTimeFormat(t *testing.T) {
	ts, err := NewTimeSheet(&http.Client{})
	if err != nil {
		t.Fatal(err)
	}
	t.Run("split time format", func(t *testing.T) {
		ts.timeDayStart, _ = time.Parse("15:04", "07:00")
		rows := []struct {
			buf             string
			day             int
			hhmmStr         string
			hhmmStrUnparsed string
		}{
			{"April 01, 2019 at 09:34AM", 1, "09:34", "09:34"},
			{"April 01, 2019 at 05:34PM", 1, "17:34", "17:34"},
			{"April 01, 2019 at 12:34PM", 1, "12:34", "12:34"},
			// midnight
			{"April 02, 2019 at 00:34AM", 1, "24:34", "00:34"},
		}

		for _, row := range rows {
			day, hhmmStr, hhmm := ts.splitTimeFormat(row.buf)
			if day != row.day {
				t.Fatalf("invalid day. actual=%d, expected=%d", day, row.day)
			}
			if hhmmStr != row.hhmmStr {
				t.Fatalf("invalid hhmmStr. actual=%s, expected=%s", hhmmStr, row.hhmmStr)
			}
			expectedHhmm, _ := time.Parse("15:04", row.hhmmStrUnparsed)
			if hhmm != expectedHhmm {
				t.Fatalf("invalid day. actual=%s, expected=%s", hhmm, expectedHhmm)
			}
		}
	})

	t.Run("time day start", func(t *testing.T) {
		strDayStart = "12:00"
		err := ts.Setup()
		if err == nil {
			t.Fatalf("timeDayStart should fail if afternoon")
		}
		t.Logf("err msg: %s", err)
	})
}
