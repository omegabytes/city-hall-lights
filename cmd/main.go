package main

import (
	"fmt"
	"os"
	"time"

	"city-hall-lights/internal/bot"
	"city-hall-lights/internal/model"
	"city-hall-lights/internal/scraper"
	"city-hall-lights/internal/store"
)

func main() {
	fs := store.NewFileStore()
	var event *model.Event
	// check if events have already been parsed into a file
	exists, err := fs.CheckFileExists()
	if err != nil {
		fmt.Println("failed to check file: ", err)
		os.Exit(1)
	}

	// if there is an event today, post it
	if exists {
		event, err = fs.Read(time.Now())
		if err != nil {
			fmt.Println("failed to read event: ", err)
		}
		if event == nil {
			fmt.Println("no event today")
			os.Exit(0)
		}
		fmt.Println(fmt.Sprintf(`today's event: %s`, event.Description))
		bot.CreateAndSendPost(event)
		os.Exit(0)
	}

	// check if events for the current month have been posted to the website
	newDataAvail, err := scraper.CheckPageLastUpdated()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !newDataAvail {
		fmt.Println("no new data available, exiting")
		os.Exit(0)
	}

	// if not, scrape the website and store the events in a file
	scrapedEvents, err := scraper.Scrape()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf(`Found %d events`, len(scrapedEvents)))
	for _, event := range scrapedEvents {
		fmt.Println(fmt.Sprintf(`%+v`, event))
	}
	if err = fs.Create(scrapedEvents); err != nil {
		fmt.Println("failed to persist events to file: ", err)
		os.Exit(1)
	}
	fmt.Println("successfully persisted events to file")
	os.Exit(0)
}

/*
Example table:

date         | color          | description                       | imageURL | defaultImageURL | altDescription                   | rawString
---------------------------------------------------------------------------------------------------------------------------------------------
11/1 - 11/6  | Red/white/blue | Election!                         |          |                 | Register to vote at: www.abc.com | "11/1 - 11/6      Red/white/blue – Election!"
11/9         | Purple 	      | Honoring our Hospitality industry |          |                 |
*/
