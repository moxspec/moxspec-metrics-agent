package main

import (
	"strings"

	"github.com/moxspec/moxspec/blk/raid"
	"github.com/moxspec/moxspec/pci"
)

func scanRAIDController(ctl pci.Device) ([]metricsReporter, error) {
	raidd := raid.NewDecoder(ctl.Path)
	err := raidd.Decode()
	if err != nil {
		return nil, err
	}

	prefix := func(t string, ts ...string) bool {
		return matcher(func(l string) bool {
			return strings.HasPrefix(ctl.Driver, l)
		}, t, ts...)
	}

	switch {
	case prefix("megaraid"):
		return scanMegaRAIDController(ctl, raidd.VirtualDisks)
	case prefix("mpt"):
		// NOTE: mpt3sas doesn't return any metrics.
	case prefix("hpvsa", "hpsa"):
		// NOTE: hpraid doesn't return any metrics.
	}
	return nil, nil
}
