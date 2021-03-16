package main

import (
	"time"

	"github.com/moxspec/moxspec/pci"
	"github.com/moxspec/moxspec/smbios"
)

func scanDevices(defInterval time.Duration) []profile {
	pcid := pci.NewDecoder()
	err := pcid.Decode()
	if err != nil {
		log.Fatal(err)
	}
	spec := smbios.NewDecoder()
	pcidevs := pci.NewDecoder()

	type Decoder interface {
		Decode() error
	}
	decoders := []Decoder{
		spec,
		pcidevs,
	}
	for _, d := range decoders {
		err := d.Decode()
		if err != nil {
			log.Warn(err.Error())
		}
	}

	// TODO: make interval configurable
	var profs []profile
	for _, p := range scanProcessor() {
		profs = append(profs, profile{
			interval: defInterval,
			reporter: p,
		})
	}

	for _, p := range scanMemory() {
		profs = append(profs, profile{
			interval: defInterval,
			reporter: p,
		})
	}

	for _, p := range scanStorage(pcidevs) {
		profs = append(profs, profile{
			interval: defInterval,
			reporter: p,
		})
	}

	for _, p := range scanNetwork(pcidevs) {
		profs = append(profs, profile{
			interval: defInterval,
			reporter: p,
		})
	}

	return profs
}
