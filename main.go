package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sclevine/agouti"
)

const (
	pathWorkTime = "/lightning/n/teamspirit__AtkWorkTimeTab"
)

var (
	domain = os.Getenv("TS_DOMAIN")

	userName = os.Getenv("TS_USER_NAME")
	password = os.Getenv("TS_PASSWORD")
)

func main() {
	if err := bulkInput(); err != nil {
		log.Fatal(err)
	}
	log.Println("success to input time sheet")
}

func bulkInput() error {
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			//TODO "--headless",
			"--disable-notifications",
		}),
	)
	defer driver.Stop()
	if err := driver.Start(); err != nil {
		return fmt.Errorf("failed to start driver: %s", err)
	}
	workTimePage, err := login(driver)
	if err != nil {
		return err
	}
	if err := focusOnTimeSheet(workTimePage); err != nil {
		return err
	}
	if err := inputWorkTime(workTimePage, 1, "10:02", false); err != nil {
		return err
	}
	return nil
}

func login(driver *agouti.WebDriver) (*agouti.Page, error) {
	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		return nil, fmt.Errorf("failed to open page: %s", err)
	}
	if err := page.Navigate(domain + pathWorkTime); err != nil {
		return nil, fmt.Errorf("failed to inputWorkTime: %s", err)
	}
	if err := page.FindByID("username").Fill(userName); err != nil {
		return nil, fmt.Errorf("failed to fill user name: %s", err)
	}
	if err := page.FindByID("password").Fill(password); err != nil {
		return nil, fmt.Errorf("failed to fill password: %s", err)
	}
	if err := page.FindByID("Login").Click(); err != nil {
		return nil, fmt.Errorf("failed to click login button: %s", err)
	}

	time.Sleep(11 * time.Second)
	return page, nil
}

func focusOnTimeSheet(page *agouti.Page) error {
	if err := page.FindByXPath(
		"//div[@class='slds-template__container']//div[@class='oneAlohaPage'][last()]//iframe",
	).SwitchToFrame(); err != nil {
		return fmt.Errorf("failed to switch to iframe: %s", err)
	}
	return nil
}

func inputWorkTime(page *agouti.Page, day int, inputTime string, isStart bool) error {
	if err := page.FindByID(fmt.Sprintf("ttvTimeSt%04d-%02d-%02d", 2019, 4, day)).Click(); err != nil {
		return fmt.Errorf("failed to click 2019/4/1: %s", err)
	}
	var index int
	if !isStart {
		index = 1
	}
	inputTag := page.FindByID("dialogInputTime").All("input").At(index)
	time.Sleep(2 * time.Second)
	if err := inputTag.Click(); err != nil {
		return fmt.Errorf("failed to click start time: %s", err)
	}
	if err := inputTag.Clear(); err != nil {
		return fmt.Errorf("failed to clear start time: %s", err)
	}
	if err := inputTag.Fill(inputTime); err != nil {
		return fmt.Errorf("failed to input start time: %s", err)
	}
	//TODO click register button
	time.Sleep(1 * time.Second)
	return nil
}
