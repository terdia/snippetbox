package main

import (
	"testing"
	"time"
)

func TestToHumanReadableDate(t *testing.T) {

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2021, 4, 17, 10, 0, 0, 0, time.UTC),
			want: "17 Apr 2021 at 10:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2021, 4, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Apr 2021 at 09:00",
		},
	}

	for _, test := range tests {
		// Use the t.Run() function to run a sub-test for each test case. The
		// first parameter to this is the name of the test (which is used to
		// identify the sub-test in any log output) and the second parameter is
		// an anonymous function containing the actual test for each case.
		t.Run(test.name, func(t *testing.T) {
			humandReadable := toHumanReadableDate(test.tm)
			if humandReadable != test.want {
				t.Errorf("want %q; got %q", test.want, humandReadable)
			}
		})
	}

}
