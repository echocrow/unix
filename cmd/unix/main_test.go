package main_test

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	unixCmd "github.com/echocrow/unix/cmd/unix"
	"github.com/stretchr/testify/assert"
)

func TestCmd(t *testing.T) {
	tests := []struct {
		args    []string
		wantOut string
	}{
		{
			[]string{"1111111111"},
			"Fri Mar 18 01:58:31 UTC 2005",
		},
		{
			[]string{"Fri Mar 18 01:58:31 UTC 2005"},
			"1111111111",
		},
		{
			[]string{"2005-03-18 01:58:31"},
			"1111111111",
		},
		{
			[]string{"1983-01-01 13:37:11", "-f", "long"},
			"Sat Jan  1 13:37:11 UTC 1983",
		},
		{
			[]string{"1983-01-01 13:37:11", "-f", "%Y-%m-%d %H:%M:%S"},
			"1983-01-01 13:37:11",
		},
		{
			[]string{"1580650631", "-f", "%Y-%m-%d %H:%M:%S"},
			"2020-02-02 13:37:11",
		},
		{
			[]string{"123456789", "-Z", "Europe/London", "-f", "%Y-%m-%d %H:%m:%m"},
			"1973-11-29 21:11:11",
		},
		{
			[]string{"2000-01-01 00:00:00", "-f", "unix"},
			"946684800",
		},
		{
			[]string{"2000-01-01 00:00:00", "-f", "unix", "-a", "13h37m11s"},
			"946733831",
		},
		{[]string{"invalid input"}, ""},
		{[]string{"1111111111", "invalid arg"}, ""},
		{[]string{"--invalid-flag"}, ""},
		{[]string{"-z", "invalid zone"}, ""},
		{[]string{"-Z", "invalid zone"}, ""},
		{[]string{"-f", "%!"}, ""},
		{[]string{"-a", "invalid offset"}, ""},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(strings.Join(tc.args, " "), func(t *testing.T) {
			t.Parallel()
			stdout, stderr, err := executeCmd(tc.args)
			wantOk := tc.wantOut != ""
			if wantOk {
				assert.Equal(t, tc.wantOut+"\n", stdout)
				assert.Empty(t, stderr)
				assert.NoError(t, err)
			} else {
				assert.NotEmpty(t, stderr)
				assert.Error(t, err)
			}
		})
	}
}

func TestNowCmd(t *testing.T) {
	tsRe := fullStdoutRe(`\d+`)
	tests := []struct {
		args []string
	}{
		{[]string{""}},
		{[]string{"now"}},
		{[]string{"now", "-z", "utc"}},
		{[]string{"-z", "utc", "-Z", "utc"}},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(strings.Join(tc.args, " "), func(t *testing.T) {
			t.Parallel()
			stdout, stderr, err := executeCmd(tc.args)
			assert.Regexp(t, tsRe, stdout)
			assert.Empty(t, stderr)
			assert.NoError(t, err)
		})
	}
}

func TestMiscCmd(t *testing.T) {
	tests := []struct {
		args     []string
		stdoutRe *regexp.Regexp
		stderrRe *regexp.Regexp
	}{
		{
			[]string{"--version"},
			fullStdoutRe(`unix version .{0,16}`),
			emptyRe,
		},
		{
			[]string{"--help"},
			regexp.MustCompile(`(?s)^Unix.+Usage:.+Flags:`),
			emptyRe,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(strings.Join(tc.args, " "), func(t *testing.T) {
			t.Parallel()
			stdout, stderr, err := executeCmd(tc.args)
			assert.Regexp(t, tc.stdoutRe, stdout)
			assert.Regexp(t, tc.stderrRe, stderr)
			assert.NoError(t, err)
		})
	}
}

func executeCmd(
	args []string,
) (stdout string, stderr string, err error) {
	cmd := unixCmd.NewCmd()
	cmd.SetArgs(args)

	bufOut := new(bytes.Buffer)
	cmd.SetOut(bufOut)

	bufErr := new(bytes.Buffer)
	cmd.SetErr(bufErr)

	err = cmd.Execute()

	stdout = bufOut.String()
	stderr = bufErr.String()
	return
}

var emptyRe = regexp.MustCompile(`^$`)

func fullStdoutRe(pattern string) *regexp.Regexp {
	return regexp.MustCompile(`^` + pattern + `\n$`)
}
