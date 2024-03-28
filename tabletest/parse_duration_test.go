//go:build !change

package tabletest

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		actual   string
		expected time.Duration
	}{
		{"0", 0},
		{"-8s", -8 * time.Second},
		{"2.4s", 2*time.Second + 400*time.Millisecond},
		{"99ns", 99 * time.Nanosecond},
		{"11m", 11 * time.Minute},
		{"-5m3.7s", -(5*time.Minute + 3*time.Second + 700*time.Millisecond)},
		{"6000000h", 0},
		{"19us", 19 * time.Microsecond},
		{"21µs", 21 * time.Microsecond},
		{"33ms", 33 * time.Millisecond},
		{"9223372036854775807ns", (1<<63 - 1) * time.Nanosecond},
		{"-9223372036854775807ns", -1<<63 + 1*time.Nanosecond},
		{"47s", 47 * time.Second},
		{"7.5h", 7*time.Hour + 30*time.Minute},
		{"2.7777s", 2*time.Second + 777700000},
		{"6h5m4s3ms2us1ns", 6*time.Hour + 5*time.Minute + 4*time.Second + 3*time.Millisecond + 2*time.Microsecond + 1*time.Nanosecond},
		{"8h30m15s", 8*time.Hour + 30*time.Minute + 15*time.Second},
		{"9223372036854775808ns", 0},
		{"9223372036854775.808us", 0},
		{"0.100000000000000000000h", 6 * time.Minute},
		{"0.830103483285477580700h", 49*time.Minute + 48*time.Second + 372539827*time.Nanosecond},
		{"", 0},
		{"9", 0},
		{"-", 0},
		{"s", 0},
		{".", 0},
		{"-.", 0},
		{".s", 0},
		{"+.s", 0},
		{"+.s", 0},
		{"+", 0},
		{"3000000h", 0},
		{"-400ms", -400 * time.Millisecond},
		{"invalid", 0},
		{"9223372036854ms775us808ns", 0},
		{"-9223372036854775808ns", 0},
	}

	for _, test := range tests {
		t.Run(test.actual, func(t *testing.T) {
			duration, err := ParseDuration(test.actual)

			if err != nil {
				if test.expected != 0 {
					t.Errorf("error in %v: %v", test.actual, err)
				}

				return
			}

			if duration != test.expected {
				t.Errorf("duration (%v): %v, expected %v", test.actual, duration, test.expected)
			}
		})
	}
}
