package images

import "fmt"

type event struct {
	Date            string
	Color           string
	Description     string
	DefaultImageURL string
	RawString       string
}

// ColorImages contains default and fetched image URLs for a given color or combination
type ColorImages struct {
	Colors          string
	DefaultImageURL string
	ImageURL        string
}

// color consts, working list. This probably needs to be converted to an appropriate data structure...
const (
	// Obvious colors
	Blue   = "Blue"
	Gold   = "Gold"
	Green  = "Green"
	Orange = "Orange"
	Purple = "Purple"
	Pink   = "Pink"
	Red    = "Red"
	Teal   = "Teal"
	White  = "White"

	// Color shades
	CarolinaBlue = "Carolina Blue"
	LightBlue    = "Light Blue"
	LimeGreen    = "Lime Green"
	SkyBlue      = "Sky Blue"

	// "mixed" colors
	Amber         = "Amber"           // thanksgiving
	Autumnal      = "Autumnal Colors" // thanksgiving
	HarvestColors = "Harvest Colors"  // thanksgiving
	Rainbow       = "Rainbow"         // pride
	RedWhitBlue   = "Red/White/Blue"  // 4th of july

	// non-standard lighting (the OOOOF zone)
	Dark        = "Dark"         // "Dark from 8:30pm to 9:30pm for WorldWide Earth Hour", "Lights will be extinguished from 8:30pm to 9:30pm for Earth Hour."
	TBD         = "TBD"          // Indiginous People's Day
	HarryPotter = "House Colors" // Harry Potter and the Cursed Child at hte Curran Theater
)

func setDefaultImageURL(event event) {
	fmt.Println(event.Color)
}

// InitSupportedColors contains a list of colors observed on sfgov.org/cityhall/lighting since 2021-01-01.
// This structure is a little weird due to the bespoke terminology used by whoever maintains the lighting schedule.
func InitSupportedColors() map[string]ColorImages {
	// We'll index on the lighting event color given to us by the maintainers of the city hall lighting schedule.
	// At first, we won't support many variations so things like "carolina blue" will be mapped to "teal" for simplicity.
	// The lighting schedule indicates that order is important - red/white/blue is not the same as white/blue/red. It reamins
	// to be seen if this has a real-world impact, so we'll ignore the distinction for now.
	supportedColors := map[string]ColorImages{
		"blue":            {"blue", "", ""},
		"gold":            {"gold", "", ""},
		"purple":          {"purple", "", ""},
		"red/green/black": {"red/green/black", "", ""},
		"red/orange":      {"red/orange", "resources/images/default-images/RedGold49ers.jpg", ""},
		"reds/orange":     {"red/orange", "resources/images/default-images/RedGold49ers.jpg", ""},
		"red/white/blue":  {"red/white/blue", "", ""},
	}

	return supportedColors
}
