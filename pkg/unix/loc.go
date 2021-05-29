package unix

import (
	"errors"
	"strings"
	"time"
)

// UNIX location errors.
var (
	ErrInvalidLocationName = errors.New("location name not recognized")
)

// Sort-opinionated time zone groups.
var tzLocGroups = []string{
	"America",
	"Europe",

	"Africa",
	"Asia",
	"Australia",
	"Brazil",
	"Canada",
	"Pacific",
	"US",
}

// ParseLocation parses a given location name and loads the time location it
// represents.
func ParseLocation(s string) (*time.Location, error) {
	if loc, err := time.LoadLocation(s); err == nil {
		return loc, nil
	}
	if loc, err := time.LoadLocation(strings.ToUpper(s)); err == nil {
		return loc, nil
	}
	if loc, err := time.LoadLocation(strings.Title(s)); err == nil {
		return loc, nil
	}

	tsLocSubS := strings.ReplaceAll(strings.Title(s), " ", "_")
	for _, tzLocGroup := range tzLocGroups {
		s := tzLocGroup + "/" + tsLocSubS
		if loc, err := time.LoadLocation(s); err == nil {
			return loc, nil
		}
	}

	return nil, ErrInvalidLocationName
}
