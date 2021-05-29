package unix_test

import (
	"testing"
	"time"

	"github.com/echocrow/unix/pkg/unix"
	"github.com/stretchr/testify/assert"
)

func TestParseLocation(t *testing.T) {
	tests := []struct {
		input string
		want  *time.Location
	}{
		{"", time.UTC},

		{"Local", time.Local},
		{"local", time.Local},

		{"UTC", time.UTC},
		{"GMT", loadTestLocation(t, "GMT")},
		{"CET", loadTestLocation(t, "CET")},
		{"EST", loadTestLocation(t, "EST")},

		{"utc", loadTestLocation(t, "utc")},

		{"Europe/Vienna", loadTestLocation(t, "Europe/Vienna")},
		{"America/Toronto", loadTestLocation(t, "America/Toronto")},
		{"America/New_York", loadTestLocation(t, "America/New_York")},

		{"vienna", loadTestLocation(t, "Europe/Vienna")},
		{"toronto", loadTestLocation(t, "America/Toronto")},
		{"new york", loadTestLocation(t, "America/New_York")},
		{"sydney", loadTestLocation(t, "Australia/Sydney")},
		{"honolulu", loadTestLocation(t, "Pacific/Honolulu")},

		{"invalid location", nil},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			got, err := unix.ParseLocation(tc.input)
			assert.Equal(t, tc.want, got)
			if tc.want != nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func loadTestLocation(t *testing.T, s string) *time.Location {
	loc, _ := time.LoadLocation(s)
	assert.NotNil(t, loc, "Failed to load test time location")
	return loc
}
