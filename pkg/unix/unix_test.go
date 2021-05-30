package unix_test

import (
	"testing"
	"time"

	"github.com/echocrow/unix/pkg/unix"
	"github.com/stretchr/testify/assert"
)

func TestParseString(t *testing.T) {
	tests := []struct {
		s      string
		wantTs int64
		wantF  string
	}{
		{"2020-01-02", 1577923200, "2006-01-02"},
		{"2020-01-02 00:00:00", 1577923200, "2006-01-02 15:04:05"},
		{"2020-01-02T00:00:00Z", 1577923200, time.RFC3339Nano},
		{"Jan 02 2020", 1577923200, "Jan 02 2006"},
		{"02 Jan 2020", 1577923200, "02 Jan 2006"},
		{"20-01-02", 1577923200, "06-01-02"},

		{"99-12-31", 946598400, "06-01-02"},
		{"31-12-99", 946598400, "02-01-06"},
		{"31.12.99", 946598400, "02.01.06"},
		{"31.12. 1999", 946598400, "02.01.2006"},
		{"31.12.1999", 946598400, "02.01.2006"},

		{"1999-12-31 23:59", 946684740, "2006-01-02 15:04"},
		{"99-12-31 23:59", 946684740, "06-01-02 15:04"},

		{"1999-12-31 23:59:58", 946684798, "2006-01-02 15:04:05"},
		{"99-12-31 23:59:58", 946684798, "06-01-02 15:04:05"},
		{"31. 12. 1999 23:59:58", 946684798, "02.01.2006 15:04:05"},
		{"31. 12. 99 23:59:58", 946684798, "02.01.06 15:04:05"},

		{"invalid input", -1, ""},
		{"99-99-99", -1, ""},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.s, func(t *testing.T) {
			t.Parallel()
			got, gotF, err := unix.Parse(tc.s)
			wantTs := tc.wantTs
			if wantTs >= 0 {
				want := time.Unix(wantTs, 0).In(time.UTC)
				assert.Equal(t, want, got)
				assert.Equal(t, tc.wantF, gotF)
				assert.NoError(t, err)
			} else {
				want := time.Time{}
				assert.Equal(t, want, got)
				assert.Error(t, err)
			}
		})
	}
}

func TestParseNow(t *testing.T) {
	tests := []struct {
		s    string
		locS string
	}{
		{"", ""},
		{"now", ""},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.s, func(t *testing.T) {
			t.Parallel()
			got, gotF, err := unix.Parse(tc.s)
			want := time.Now().Round(0)
			accuracy := time.Second
			minWant := want.Add(accuracy * -1)
			maxWant := want.Add(accuracy)
			assert.True(t, minWant.Before(got), "want time >= min threshold")
			assert.True(t, got.Before(maxWant), "want time <= max threshold")
			assert.Equal(t, "now", gotF)
			assert.NoError(t, err)
		})
	}
}
