package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCreateMetricsList(t *testing.T) {
	var err error
	var got []metrics
	tests := []struct {
		labels     map[string]string
		points     map[string]interface{}
		ex         []metrics
		wantsError bool
	}{
		{nil, nil, nil, true},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got, err = createMetricsList(tt.labels, tt.points)
			if err == nil && tt.wantsError || err != nil && !tt.wantsError || !reflect.DeepEqual(got, tt.ex) {
				t.Errorf("test: %+v, got: %v (err: %s), expect: %v (wantsError: %t)", tt, got, err, tt.ex, tt.wantsError)
			}
		})
	}
}

func TestNewMetrics(t *testing.T) {
	var err error
	var got metrics

	tLabel := map[string]string{
		"kind": "test",
	}
	tests := []struct {
		name       string
		labels     map[string]string
		val        interface{}
		ex         metrics
		wantsError bool
	}{
		{"", nil, 0, metrics{}, true},
		{"test", nil, 0, metrics{}, true},
		{"test", tLabel, 1.0, metrics{Name: getMetricsNameWithPrefix("test"), Labels: tLabel, Value: 1.0}, false},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got, err = newMetrics(tt.name, tt.labels, tt.val)
			isSameVal := (got.Name == tt.ex.Name && reflect.DeepEqual(got.Labels, tt.ex.Labels) && got.Value == tt.ex.Value)
			if err == nil && tt.wantsError || err != nil && !tt.wantsError || !isSameVal {
				t.Errorf("test: %+v, got: %v (err: %s), expect: %v (wantsError: %t)", tt, got, err, tt.ex, tt.wantsError)
			}
		})
	}

}

func TestMetricsNameGenerator(t *testing.T) {
	var got string
	tests := []struct {
		name string
		ex   string
	}{
		{"248", prefix + "_248"},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = getMetricsNameWithPrefix(tt.name)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %s, expect: %s", tt, got, tt.ex)
			}
		})
	}

}

func TestCopyLabels(t *testing.T) {
	var got map[string]string
	tests := []struct {
		in map[string]string
		ex map[string]string
	}{
		{map[string]string{"248": "248"}, map[string]string{"248": "248"}},
		{map[string]string{}, map[string]string{}},
		{nil, nil},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = copyLabels(tt.in)
			if !reflect.DeepEqual(got, tt.ex) {
				t.Errorf("test: %+v, got: %+v, expect: %+v", tt, got, tt.ex)
			}
		})
	}

}
