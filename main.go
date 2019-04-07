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
	if err := inputWorkTime(workTimePage); err != nil {
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

	time.Sleep(13 * time.Second)
	return page, nil
}

func inputWorkTime(page *agouti.Page) error {
	//TODO xpath page.FindByXPath("//div[@class='slds-template__container'][-1]//iframe")
	selector := page.FindByClass("slds-template__container").AllByClass("oneAlohaPage").At(-1)
	if err := selector.FindByXPath("//iframe").SwitchToFrame(); err != nil {
		return fmt.Errorf("failed to switch to iframe: %s", err)
	}
	if err := page.FindByID("ttvTimeSt2019-04-01").Click(); err != nil {
		//if err := page.FindByID(fmt.Sprintf("ttvTimeSt%04d-%02d-%02d", 2019, 4, 1)).Click(); err != nil {
		return fmt.Errorf("failed to click 2019/4/1: %s", err)
	}
	inputTag := page.FindByID("dialogInputTime").All("input").At(0)
	time.Sleep(2 * time.Second)
	if err := inputTag.Click(); err != nil {
		return fmt.Errorf("failed to click start time: %s", err)
	}
	if err := inputTag.Clear(); err != nil {
		return fmt.Errorf("failed to clear start time: %s", err)
	}
	if err := inputTag.Fill("10:01"); err != nil {
		return fmt.Errorf("failed to input start time: %s", err)
	}
	//TODO click register button
	time.Sleep(1 * time.Second)
	return nil
}
