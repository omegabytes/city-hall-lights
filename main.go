package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type event struct {
	Date            string
	Color           string
	Description     string
	DefaultImageURL string
	RawString       string
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("sfgov.org"),
	)

	events := make([]event, 0)
	c.OnHTML(".field-item", func(e *colly.HTMLElement) {
		fmt.Println("found field-item")
		rawEventString := e.ChildText("p")

		events = parseEvents(rawEventString)
	})

	c.OnRequest(func(r *colly.Request) {
		// todo: check for last updated, if possible
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://sfgov.org/cityhall/lighting")

	if len(events) > 0 {
		for _, event := range events {
			printMessage(event)
			if isEventToday(event.Date) {
				printMessage(event)
			}
		}
	} else {
		fmt.Println(fmt.Errorf("lighting schedule list was empty"))
	}

}

// parseEvents scrapes "https://sfgov.org/cityhall/lighting" for the html element containing a
// schedule of lighting events, each with a color or description. This is extremly fragile and will
// only work for this url.
// TODO: make generic for use on https://sfgov.org/cityhall/past-lighting-schedule-2020
func parseEvents(rawEventString string) []event {
	dateRegex, _ := regexp.Compile("[^A-Za-z]*")
	colorRegex, _ := regexp.Compile("[^-–]*")
	descRegex, _ := regexp.Compile("[^A-Za-z0-9]*")

	events := make([]event, 0)
	rawArr := strings.Split(rawEventString, "\n")

	for _, s := range rawArr {
		event := &event{}
		event.RawString = s

		date := dateRegex.FindString(s)
		// Exclude empty lines and unrelated text included in our scraped <p>
		if len(date) <= 1 {
			continue
		}
		event.Date = strings.TrimSpace(date)
		s = strings.ReplaceAll(s, date, "")

		// todo: get and format color function
		event.Color = strings.TrimRight(colorRegex.FindString(s), " ")
		s := strings.ReplaceAll(s, event.Color, "")

		// todo: get and format description function
		descDelim := descRegex.FindString(s)
		event.Description = strings.TrimRight(strings.ReplaceAll(s, descDelim, ""), " ")

		events = append(events, *event)
	}
	return events
}

// parseDate returns the integer value of the month and day of an event date.
// We expect a format of "m/d" (eg. 1/11 for Jan 11th).
// TODO: support for date ranges (eg. "1/11 - 1/13")
func parseDate(date string) (month int, day int) {
	dates := strings.Split(date, "/")
	if len(dates) > 0 {
		m, _ := strconv.Atoi(dates[0])
		d, _ := strconv.Atoi(dates[1])
		return m, d
	} else {
		fmt.Println(fmt.Errorf("couldn't parse eventDate"))
	}
	return 0, 0
}

// isEventToday returns true if the event date matches todays date. We only care about month and day.
func isEventToday(eventDate string) bool {
	today := time.Now()
	_, tm, td := today.Date()

	em, ed := parseDate(eventDate)
	if int(tm) == em && td == ed {
		return true
	}
	return false
}

func printMessage(event event) {
	fmt.Println(event.Date, event.Color, event.Description)
}

/*
Example table:

date         | color          | description                       | imageURL | defaultImageURL | altDescription                   | rawString
---------------------------------------------------------------------------------------------------------------------------------------------
11/1 - 11/6  | Red/white/blue | Election!                         |          |                 | Register to vote at: www.abc.com | "11/1 - 11/6      Red/white/blue – Election!"
11/9         | Purple 	      | Honoring our Hospitality industry |          |                 |


Strings with pattern:
11/1 - 11/6      Red/white/blue – Election!\n

raw = "11/1 - 11/6      Red/white/blue – Election!"

Parsing actions:
1. Find dates
    Numeric
    Can be range
    Can be single date
    Contains "/", " ", and "-"
	Whitespace terminated

	[^A-Za-z]* // match until first "word" (non digit, non space, non special char)
	dates = "11/1 - 11/6      "
	raw = "Red/white/blue – Election!"


2. Find colors
    Contains "/", " "
    Sometimes is defined colors (red, green)
	Sometimes is subjective ("Autumnal")
	Whitespace terminated OR "-" terminated

	[^-–]* // search until the first occurence of either en dash(U+2013) or hyphen minus (U+002D)
	color = "Red/white/blue "
	raw = "– Election!"

3. Find description
    Always* separated from colors by "-"
	End of string, potentially trailing whitespace
	Can contain special chars


	[^A-Za-z0-9]* // matches until the first alphanumeric char
	desc = "2020 Goldman Environmental Prize (31st Annual)"


Quirks
- two "-" chars are used (inconsistently)
	EN DASH – Unicode: U+2013
	HYPHEN MINUS - Unicode: U+002D

Todo dump
- validate date format (currently m/dd but could be m/dd or m/dd/yyyy)
*/
