package scraper

import (
	"fmt"
	"strings"
	"time"

	"city-hall-lights/internal/model"
	"city-hall-lights/internal/parser"
	"github.com/gocolly/colly/v2"
)

const (
	lightingScheduleURL        = "https://www.sf.gov/location/san-francisco-city-hall"
	excludeFirstListElemString = "City Hall will be lit"
	excludeLastListElemString  = "Learn more about City Hall's exterior lighting and see past lighting schedules."
	selector                   = `#block-sfgovpl-content > article > div.sfgov-section-container > 
								  div.group--left > div.sfgov-section.sfgov-section-getting-here > div > 
								  div.field.field--type-entity-reference-revisions.__getting-here-items.field__items > 
								  div:nth-child(5) > details > div > div`
)

func Scrape() ([]model.Event, error) {
	events := make([]model.Event, 0)

	c := colly.NewCollector()

	c.OnHTML(selector, func(e *colly.HTMLElement) {
		e.ForEach("p", func(_ int, el *colly.HTMLElement) {
			parsedEvent := model.Event{}
			if el.Text != "\u00a0" && !strings.Contains(el.Text, excludeFirstListElemString) &&
				!strings.Contains(el.Text, excludeLastListElemString) {
				parsedEvent = parser.ParseEvent(el.Text)
				events = append(events, parsedEvent)
			}
		})
	})

	if err := c.Visit(lightingScheduleURL); err != nil {
		return nil, err
	}
	return events, nil
}

func CheckPageLastUpdated() (bool, error) {
	newDataAvailable := false

	c := colly.NewCollector()
	lastUpdatedSelector := "#block-sfgovpl-content > article > div.sfgov-section-container > div.group--left > div.sfgov-section.sfgov-section-getting-here > div > div.field.field--type-entity-reference-revisions.__getting-here-items.field__items > div:nth-child(5) > details > div > div > h4"
	c.OnHTML(lastUpdatedSelector, func(e *colly.HTMLElement) {
		fmt.Println("date: ", e.Text)
		availEvents, err := ParseAvailableEventsMonthYear(e.Text)
		if err != nil {
			fmt.Println(err)
		}
		thisMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
		fmt.Println("this month: ", thisMonth)
		fmt.Println("events published for: ", availEvents)
		if thisMonth.Equal(availEvents) {
			fmt.Println("new data available")
			newDataAvailable = true
			return
		}
		fmt.Println("no new data")
	})

	if err := c.Visit(lightingScheduleURL); err != nil {
		return newDataAvailable, err
	}
	return true, nil
}

func ParseAvailableEventsMonthYear(rawString string) (time.Time, error) {
	layout := "January 2006"
	parts := strings.Split(rawString, " ")
	return time.Parse(layout, fmt.Sprintf("%s %s", parts[0], parts[1]))
}
