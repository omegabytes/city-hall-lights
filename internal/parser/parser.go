package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"city-hall-lights/internal/model"
)

// colorSplitRE is used to split the color from the raw event string. Characters used are unfortunately not uniform,.
var colorSplitRE = regexp.MustCompile("[-–]")

// ParseEvent attempts to parse a raw event string into a structured event.
// Since the data is not consistent, this method has a fair amount of extra complexity.
func ParseEvent(rawEventString string) model.Event {
	rawEventString = strings.Replace(rawEventString, ` `, ``, -1)
	parts := colorSplitRE.Split(rawEventString, -1)

	if len(parts) == 3 {
		date := strings.TrimRight(parts[0], " ")
		startTimeStamp, err := convertToTimestamps(date)
		if err != nil {
			fmt.Println(err)
			return model.Event{
				RawEventString: rawEventString,
			}
		}
		color := strings.TrimSpace(parts[1])
		description := transformDescription(parts[2], color)
		return model.Event{
			DateString:     date,
			StartTimeStamp: startTimeStamp[0],
			Color:          color,
			Description:    description,
			RawEventString: rawEventString,
		}
	}

	return model.Event{
		RawEventString: rawEventString,
	}
}

// convertToTimestamps attempts to convert a date string into a timestamp.
// The language used to describe the date is not consistent.
// The returned timestamp is needed for use in automated posting on a schedule.
func convertToTimestamps(dateString string) ([]time.Time, error) {
	layout := "Monday, January 2, 2006"      // Standard format with year
	layoutWithoutYear := "Monday, January 2" // Format without year
	year := 2024                             // Set the target year here, if known

	// Trim and check if input contains "and" or "through"
	dateString = strings.TrimSpace(dateString)
	if strings.Contains(dateString, " and ") {
		// Split and parse each date individually
		parts := strings.Split(dateString, " and ")
		dates := []time.Time{}

		for _, part := range parts {
			date, err := parseSingleDate(part, layout, layoutWithoutYear, year)
			if err != nil {
				return nil, err
			}
			dates = append(dates, date)
		}
		return dates, nil
	}

	if strings.Contains(dateString, " through ") {
		// Split and parse start and end dates for a range
		parts := strings.Split(dateString, " through ")
		if len(parts) != 2 {
			return nil, errors.New("invalid range format")
		}

		startDate, err := parseSingleDate(parts[0], layout, layoutWithoutYear, year)
		if err != nil {
			return nil, err
		}
		endDate, err := parseSingleDate(parts[1], layout, layoutWithoutYear, year)
		if err != nil {
			return nil, err
		}

		return []time.Time{startDate, endDate}, nil
	}

	// If it's a single date
	date, err := parseSingleDate(dateString, layout, layoutWithoutYear, year)
	if err != nil {
		return nil, err
	}
	return []time.Time{date}, nil
}

func parseSingleDate(dateStr, layout, layoutWithoutYear string, year int) (time.Time, error) {
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load location: %w", err)
	}
	date, err := time.ParseInLocation(layout, dateStr, location)
	if err != nil {
		// Try parsing without the year and add it manually
		date, err = time.ParseInLocation(layoutWithoutYear, dateStr, location)
		if err != nil {
			return time.Time{}, err
		}
		date = date.AddDate(year, 0, 0) // Add the target year
	}
	return date, nil
}

const (
	recognitionPattern = "in recognition of "
	commemoratePattern = "to commemorate "
)

func transformDescription(rawDescription string, color string) string {
	colors, _ := transformColors(color)
	desc := strings.TrimLeft(rawDescription, " ")
	desc = strings.Replace(desc, `“`, `"`, -1)

	joiner := ""
	if !strings.Contains(desc, recognitionPattern) && !strings.Contains(desc, commemoratePattern) {
		joiner = recognitionPattern
	}
	descTemplate := fmt.Sprintf(`Tonight City Hall will be %s %s%s`, colors, joiner, desc)
	return strings.TrimSpace(descTemplate)
}

// transformColors converts input strings into a formatted list. Leading and trailing spaces are trimmed,
// and input is lowercased. Oxford comma style is used for three or more colors.
// Known formats:
// * "red" -> "red"
// * "purple/yellow" -> "purple and yellow"
// * "red/white/blue" -> "red, white, and blue"
// * "shades of amber" -> "shades of amber"
func transformColors(colors string) (string, error) {
	if strings.TrimSpace(colors) == "" {
		return "", fmt.Errorf("invalid input: color must not be empty")
	}
	loweredInput := strings.ToLower(colors)
	parts := strings.Split(loweredInput, "/")

	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	switch len(parts) {
	case 1:
		// Single color, return as is.
		return parts[0], nil
	case 2:
		// Two colors, join with "and".
		return fmt.Sprintf("%s and %s", parts[0], parts[1]), nil
	case 3:
		// Three colors, Oxford comma style.
		return fmt.Sprintf("%s, %s, and %s", parts[0], parts[1], parts[2]), nil
	default:
		// Handle any number of parts greater than 3.
		last := parts[len(parts)-1]
		return fmt.Sprintf("%s, and %s", strings.Join(parts[:len(parts)-1], ", "), last), nil
	}
}
