package main

import (
	"testing"
	"time"
)

var tests = []struct {
	name string
	tm   time.Time
	want string
}{
	{
		name: "UTC",
		tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC),
		want: "17 Mar 2024 at 10:15",
	},
	{
		name: "Empty",
		tm:   time.Time{},
		want: "",
	},
	{
		name: "CET",
		tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
		want: "17 Mar 2024 at 09:15",
	},
}

func TestHumanDate(t *testing.T) {
	// Initialize a new time.Time value and pass it to the humanDate function.
	for _, tt := range tests {
		// Use the t.Run() function to run a sub-test for each test case. The
		// first parameter to this is the name of the test (used to identify the
		// sub-test in any log output) and the second parameter is an anonymous
		// function containing the actual test for each case.
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			if hd != tt.want {
				// Check that the output from the humanDate function is in the format we
				// expect. If it isn't what we expect, use the t.Errorf() function to
				// indicate that the test has failed and log the expected and actual
				// values.
				t.Errorf("got %q; want %q", hd, tt.want)
			}
		})
	}
}
