package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp/runner"

	cdp "github.com/chromedp/chromedp"
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
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := cdp.New(ctxt, cdp.WithRunnerOptions(
		runner.Flag("enable-automation", true),
		runner.Flag("disable-notifications", true),
	))
	if err != nil {
		log.Fatal(err)
	}

	// login
	if err := c.Run(ctxt, login()); err != nil {
		log.Fatalf("failed to login: %s", err)
	}

	// shutdown chrome
	if err := c.Shutdown(ctxt); err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	if err := c.Wait(); err != nil {
		log.Fatal(err)
	}
}

func login() cdp.Tasks {
	return cdp.Tasks{
		cdp.Navigate(domain),
		//cdp.WaitVisible("#username"),
		cdp.SendKeys("//input[@id='username']", userName, cdp.NodeVisible),
		cdp.Click("#password"),
		cdp.SendKeys("//input[@id='password']", password),
		cdp.Click("#Login"),
		cdp.Sleep(2 * time.Second), //FIXME
		cdp.Navigate(domain + pathWorkTime),
		cdp.Sleep(10 * time.Second),
		cdp.Click("#holyLink"),
	}
}
