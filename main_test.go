package main

import (
	"fmt"
	"testing"
	"time"
)

func TestDryrun(t *testing.T) {
	dummyRep := func() ([]metrics, error) {
		return []metrics{
			{
				name:      "dummy",
				labels:    map[string]string{"dummy": "dummy"},
				timestamp: time.Now(),
				value:     0.1,
			},
		}, nil
	}

	dummyRepErr := func() ([]metrics, error) {
		return nil, fmt.Errorf("dummy err")
	}

	tests := []struct {
		in         []profile
		exLenProfs int
		exMps      float64
		wantsError bool
	}{
		{nil, 0, 0, true},
		{[]profile{{dummyRep, time.Second}}, 1, 1.0, false},
		{[]profile{{dummyRepErr, time.Second}}, 0, 0, true},
		{[]profile{{dummyRep, 0}}, 0, 0, true},
		{[]profile{{dummyRep, time.Second}, {dummyRep, time.Second / 5}}, 2, 6.0, false},
		{[]profile{{dummyRep, 0}, {dummyRep, time.Second / 5}}, 1, 5.0, false},
		{[]profile{{dummyRep, time.Second}, {dummyRep, time.Second / 5}, {dummyRepErr, time.Second}}, 2, 6.0, false},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			gotProfs, gotMps, err := dryrun(tt.in)
			if err == nil && tt.wantsError || err != nil && !tt.wantsError || len(gotProfs) != tt.exLenProfs {
				t.Errorf("test: %+v, got: profs: %+v, mps: %.1f (err: %s), expect: lenProfs=%d, mps: %.1f(wantsError: %t)",
					tt, gotProfs, gotMps, err, tt.exLenProfs, tt.exMps, tt.wantsError,
				)
			}
		})
	}
}

func TestCalcBufferSize(t *testing.T) {
	tests := []struct {
		interval   time.Duration
		mps        float64
		ex         int
		wantsError bool
	}{
		{0, 0, 0, true},
		{0, 1, 0, true},
		{time.Second, 1, 1 * bufMargin, false},
		{time.Second / 2, float64(bufMax * 2), bufMax, false},
		{time.Second * 10, 6.4, 140, false},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got, err := calcBufferSize(tt.interval, tt.mps)
			if err == nil && tt.wantsError || err != nil && !tt.wantsError || got != tt.ex {
				t.Errorf("test: %+v, got: %d (err: %s), expect: %d (wantsError: %t)", tt, got, err, tt.ex, tt.wantsError)
			}
		})
	}
}
