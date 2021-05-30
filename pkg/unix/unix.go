// Package unix provides UNIX helper functions.
package unix

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// UNIX parse errors.
var (
	ErrInvalidTimeLayout = errors.New("date input format not recognized")
)

// Special UNIX parse layouts.
const (
	NowLayout = "now"
)

// Standard time layouts, based on Go's layout template:
// Mon Jan 2 15:04:05 -0700 MST 2006
var stdTimeLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	time.StampNano,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.ANSIC,
	time.UnixDate,
	time.RFC822,
	time.RFC822Z,
}

// Allowed date layouts.
var dateLayouts = []string{
	"2006-01-02",
	"2006/01/02",
	"06-01-02",
	"02-01-06",
	"02.01.2006",
	"02.01.06",
	"02/01/06",
	"01/02/06",
	"Jan 02 2006",
	"02 Jan 2006",
}

// Allowed time layouts
var timeLayout = []string{
	"15:04:05",
	"15:04",
	"",
}

// Parse auto-detects the date/time format of a string and parses it to the
// time value it represents.
func Parse(
	s string,
	locS string,
) (t time.Time, layout string, err error) {
	if s == "" || s == NowLayout {
		var loc *time.Location
		loc, err = ParseLocation(locS)
		if err != nil {
			return
		}
		return time.Now().In(loc), NowLayout, nil
	}

	for _, layout = range stdTimeLayouts {
		if t, err = parseInLoc(layout, s, locS); err == nil {
			return
		}
	}

	for _, dateL := range dateLayouts {
		for _, timeL := range timeLayout {
			if t, layout, err = parseDateTime(dateL, timeL, s, locS); err == nil {
				return
			}
		}
	}

	ss := shortenStr(s)
	for _, dateL := range dateLayouts {
		dateL = shortenStr(dateL)
		for _, timeL := range timeLayout {
			if t, layout, err = parseDateTime(dateL, timeL, ss, locS); err == nil {
				return
			}
		}
	}

	return time.Time{}, "", ErrInvalidTimeLayout
}

func shortenStr(s string) string {
	s = regexp.MustCompile(`\s{2,}`).ReplaceAllString(s, "")
	s = regexp.MustCompile(`\s(\W)`).ReplaceAllString(s, "$1")
	s = regexp.MustCompile(`(\W)\s`).ReplaceAllString(s, "$1")
	return s
}

func parseDateTime(
	dateL string,
	timeL string,
	input string,
	locS string,
) (t time.Time, layout string, err error) {
	layout = strings.Trim(dateL+" "+timeL, " ")
	if t, err = parseInLoc(layout, input, locS); err == nil {
		return
	}

	layout = strings.Trim(timeL+" "+dateL, " ")
	if t, err = parseInLoc(layout, input, locS); err == nil {
		return
	}

	return time.Time{}, "", ErrInvalidTimeLayout
}

func parseInLoc(
	f string,
	s string,
	locS string,
) (time.Time, error) {
	if locS == "" {
		return time.Parse(f, s)
	}

	loc, err := ParseLocation(locS)
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(f, s, loc)
}
