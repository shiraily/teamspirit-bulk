package main

import (
	"context"
	"log"
	"os"

	"github.com/chromedp/chromedp"
)

var (
	domain   = os.Getenv("TS_DOMAIN")
	userName = os.Getenv("TS_USER_NAME")
	password = os.Getenv("TS_PASSWORD")
)

func main() {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// login
	err = c.Run(ctxt, login())
	if err != nil {
		log.Fatal(err)
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

func login() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(domain),
		chromedp.WaitVisible("#username"),
		chromedp.SendKeys("//input[@id='username']", userName),
		chromedp.Click("#password", chromedp.NodeVisible),
		chromedp.SendKeys("//input[@id='password']", password),
		chromedp.Click("#Login", chromedp.NodeVisible),
		//chromedp.WaitVisible("#mainTableBody"),
		//chromedp.Sleep(2 * time.Second),
	}
}
