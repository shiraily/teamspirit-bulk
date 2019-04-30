package teamspirit

import (
	"fmt"
	"log"
	"time"

	"github.com/sclevine/agouti"

	"github.com/shiraily/teamspirit-bulk/config"
	"github.com/shiraily/teamspirit-bulk/model"
)

const (
	pathWorkTime = "/lightning/n/teamspirit__AtkWorkTimeTab"
)

var (
	DefaultDriver = agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			//"--headless",
			"--disable-notifications",
		}),
	)
)

//TODO defer Driver.Stop
type TeamSpirit struct {
	Driver *agouti.WebDriver
	page   *agouti.Page

	year  int
	month time.Month
	day   int
}

func NewTeamSpirit(driver *agouti.WebDriver) *TeamSpirit {
	year, month, day := time.Now().Date()
	return &TeamSpirit{
		Driver: driver,
		year:   year,
		month:  month,
		day:    day,
	}
}

// TODO: close driver when some error
func (t *TeamSpirit) Setup() error {
	if err := t.Driver.Start(); err != nil {
		return fmt.Errorf("failed to start Driver: %s", err)
	}
	var err error
	t.page, err = t.login()
	if err != nil {
		return err
	}
	if err := t.focusOnTimeSheet(); err != nil {
		return err
	}
	return nil
}

func (t *TeamSpirit) login() (*agouti.Page, error) {
	page, err := t.Driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		return nil, fmt.Errorf("failed to open page: %s", err)
	}
	if err := page.Navigate(config.Cfg.TeamSpirit.Domain + pathWorkTime); err != nil {
		return nil, fmt.Errorf("failed to inputWorkTime: %s", err)
	}
	if err := page.FindByID("username").Fill(config.Cfg.TeamSpirit.User); err != nil {
		return nil, fmt.Errorf("failed to fill user name: %s", err)
	}
	if err := page.FindByID("password").Fill(config.Cfg.TeamSpirit.Password); err != nil {
		return nil, fmt.Errorf("failed to fill password: %s", err)
	}
	if err := page.FindByID("Login").Click(); err != nil {
		return nil, fmt.Errorf("failed to click login button: %s", err)
	}
	return page, nil
}

func (t *TeamSpirit) focusOnTimeSheet() error {
	return retry(func() error {
		if errFind := t.page.FindByXPath(
			"//div[@class='slds-template__container']//div[@class='oneAlohaPage'][last()]//iframe",
		).SwitchToFrame(); errFind != nil {
			return fmt.Errorf("failed to switch to iframe: %s", errFind)
		}
		return nil
	}, 20)
}

func (t *TeamSpirit) BulkInput(workTimes []model.WorkTime) error {
	var err error
	for _, workTime := range workTimes {
		if err := t.Input(workTime); err != nil {
			log.Print(err)
			break
		}
	}
	return err
}

func (t *TeamSpirit) Input(workTime model.WorkTime) error {
	if err := retry(func() error {
		if errFind := t.page.FindByID(
			fmt.Sprintf("ttvTimeSt%04d-%02d-%02d", t.year, t.month, workTime.Day),
		).Click(); errFind != nil {
			return fmt.Errorf("failed to click %04d/%02d/%02d: %s", t.year, t.month, workTime.Day, errFind)
		}
		return nil
	}, 10); err != nil {
		return err
	}
	dialog := t.page.FindByID("dialogInputTime")
	if err := inputTime(dialog, "startTime", workTime.StartTime); err != nil {
		return err
	}
	if err := inputTime(dialog, "endTime", workTime.EndTime); err != nil {
		return err
	}
	if err := dialog.FindByID("dlgInpTimeOk").Click(); err != nil {
		return fmt.Errorf("failed to click OK button: %s", err)
	}
	return nil
}

func inputTime(dialog *agouti.Selection, tagName, inputTime string) error {
	inputTag := dialog.FindByID(tagName)
	if err := inputTag.Click(); err != nil {
		return fmt.Errorf("failed to click %s: %s", tagName, err)
	}
	if err := inputTag.Clear(); err != nil {
		return fmt.Errorf("failed to clear %s: %s", tagName, err)
	}
	if err := inputTag.Fill(inputTime); err != nil {
		return fmt.Errorf("failed to input %s: %s", tagName, err)
	}
	return nil
}

func retry(fn func() error, count int) error {
	var err error
	for i := 0; i < count; i++ {
		err = fn()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	return err
}
