package main

import (
	"strings"

	"github.com/moxspec/moxspec/pci"
)

// for switch-selecter
func matcher(f func(string) bool, t string, ts ...string) bool {
	list := append([]string{t}, ts...)
	for _, l := range list {
		if f(l) {
			return true
		}
	}
	return false
}

func scanStorage(pcidevs *pci.Devices) []metricsReporter {
	var reporters []metricsReporter
	for _, ctl := range pcidevs.FilterByClass(pci.MassStorageController) {
		equal := func(t string, ts ...string) bool {
			return matcher(func(l string) bool {
				return (ctl.Driver == l)
			}, t, ts...)
		}

		prefix := func(t string, ts ...string) bool {
			return matcher(func(l string) bool {
				return strings.HasPrefix(ctl.Driver, l)
			}, t, ts...)
		}

		var err error
		switch {
		case equal("nvme"):
			rs, err := scanNVMeController(ctl.Path, ctl.Driver)
			if err == nil {
				reporters = append(reporters, rs...)
			}
		case equal("ahci", "ata_piix", "isci"):
			rs, err := scanAHCIController(ctl.Path, ctl.Driver)
			if err == nil {
				reporters = append(reporters, rs...)
			}
		case prefix("mpt", "megaraid", "hpvsa", "hpsa"):
			rs, err := scanRAIDController(*ctl)
			if err == nil {
				reporters = append(reporters, rs...)
			}
		case equal("virtio-pci"):
		}

		if err != nil {
			log.Warn(err.Error())
		}

	}
	return reporters
}
