package unix

import (
	"errors"
	"fmt"
	"time"
)

// Special time format layouts.
const (
	UnixFormat = "unix"
	LongFormat = "long"
)

// Format errors.
var (
	ErrUnsupportedFormatDirective = errors.New("format directive is invalid or unsupported")
)

// Format returns a string representation of time t according to layout.
func Format(t time.Time, layout string) (string, error) {
	if layout == UnixFormat {
		return fmt.Sprint(t.Unix()), nil
	}
	if layout == LongFormat {
		layout = time.UnixDate
	} else {
		var err error
		if layout, err = decodeStandardDateDirectives(layout); err != nil {
			return "", err
		}
	}

	return t.Format(layout), nil
}

var standardDateDirectives = map[rune]string{
	'a': "Mon",
	'A': "Monday",
	'b': "Jan",
	'B': "January",
	'c': time.ANSIC,
	'd': "02",
	'f': "000000",
	'G': "2006",
	'H': "15",
	'I': "03",
	'm': "01",
	'M': "04",
	'p': "PM",
	'S': "05",
	'x': "01/02/06",
	'X': "15:04:05",
	'y': "06",
	'Y': "2006",
	'z': "-0700",
	'Z': "MST",
	'%': "%",
}

func decodeStandardDateDirectives(format string) (layout string, err error) {
	isPrefixed := false
	for _, char := range format {
		isPrefix := char == '%' && !isPrefixed
		if isPrefixed {
			goCode, ok := standardDateDirectives[char]
			if !ok {
				return "", ErrUnsupportedFormatDirective
			}
			layout += goCode
		} else if !isPrefix {
			layout += string(char)
		}
		isPrefixed = isPrefix
	}

	return layout, nil
}
