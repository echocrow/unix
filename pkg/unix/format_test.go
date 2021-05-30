package unix_test

import (
	"testing"
	"time"

	"github.com/echocrow/unix/pkg/unix"
	"github.com/stretchr/testify/assert"
)

var (
	testT0 = time.Date(2008, 9, 17, 20, 4, 26, 0, time.UTC)
	testT1 = time.Date(1994, 9, 17, 20, 4, 26, 0, time.FixedZone("EST", -18000))
	testT2 = time.Date(2000, 12, 26, 1, 15, 6, 0, time.FixedZone("OTO", 15600))
)

func TestFormat(t *testing.T) {
	tests := []struct {
		t      time.Time
		layout string
		want   string
	}{
		{testT0, "unix", "1221681866"},
		{testT1, "unix", "779850266"},
		{testT2, "unix", "977777706"},

		{time.Unix(123456789, 0), "unix", "123456789"},
		{time.Unix(987456, 0), "unix", "987456"},

		{testT0, "long", "Wed Sep 17 20:04:26 UTC 2008"},
		{testT1, "long", "Sat Sep 17 20:04:26 EST 1994"},
		{testT2, "long", "Tue Dec 26 01:15:06 OTO 2000"},

		{testT0, time.UnixDate, "Wed Sep 17 20:04:26 UTC 2008"},
		{testT1, time.UnixDate, "Sat Sep 17 20:04:26 EST 1994"},
		{testT2, time.UnixDate, "Tue Dec 26 01:15:06 OTO 2000"},

		{testT0, time.RFC3339, "2008-09-17T20:04:26Z"},
		{testT1, time.RFC3339, "1994-09-17T20:04:26-05:00"},
		{testT2, time.RFC3339, "2000-12-26T01:15:06+04:20"},

		{testT0, time.Stamp, "Sep 17 20:04:26"},
		{testT1, time.Stamp, "Sep 17 20:04:26"},
		{testT2, time.Stamp, "Dec 26 01:15:06"},

		{testT0, "dateless format", "dateless format"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.t.String()+" "+tc.layout, func(t *testing.T) {
			t.Parallel()
			got, err := unix.Format(tc.t, tc.layout)
			assertFormat(t, tc.want, got, err)
		})
	}
}

func TestFormatStandardDirectives(t *testing.T) {
	tests := []struct {
		t      time.Time
		layout string
		want   string
	}{
		{testT0, "%Y-%m-%d %H:%M:%S", "2008-09-17 20:04:26"},
		{testT1, "%Y-%m-%d %H:%M:%S", "1994-09-17 20:04:26"},
		{testT2, "%Y-%m-%d %H:%M:%S", "2000-12-26 01:15:06"},

		{testT0, "%a %b %d %H:%M:%S %Z %Y", "Wed Sep 17 20:04:26 UTC 2008"},
		{testT1, "%a %b %d %H:%M:%S %Z %Y", "Sat Sep 17 20:04:26 EST 1994"},
		{testT2, "%a %b %d %H:%M:%S %Z %Y", "Tue Dec 26 01:15:06 OTO 2000"},

		{testT0, "%c", "Wed Sep 17 20:04:26 2008"},
		{testT1, "%c", "Sat Sep 17 20:04:26 1994"},
		{testT2, "%c", "Tue Dec 26 01:15:06 2000"},

		{testT0, "%Y %%m %%%d", "2008 %m %17"},
		{testT1, "%%%Y%m%%d", "%199409%d"},
		{testT2, "%%%Y.%m.%%d", "%2000.12.%d"},

		{testT0, "%?", ""},
		{testT1, "%.", ""},
		{testT2, "%-", ""},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.t.String()+" "+tc.layout, func(t *testing.T) {
			t.Parallel()
			got, err := unix.Format(tc.t, tc.layout)
			assertFormat(t, tc.want, got, err)
		})
	}
}

func assertFormat(t *testing.T, want, got string, err error) {
	wantOk := want != ""
	assert.Equal(t, want, got)
	if wantOk {
		assert.NoError(t, err)
	} else {
		assert.Error(t, err)
	}
}
