package main

import (
	"github.com/moxspec/moxspec/netlink"
	"github.com/moxspec/moxspec/nw"
	"github.com/moxspec/moxspec/pci"
)

func scanNetwork(pcidevs *pci.Devices) []metricsReporter {
	var reporters []metricsReporter
	for _, ctl := range pcidevs.FilterByClass(pci.NetworkController) {
		nwd := nw.NewDecoder(ctl.Path, ctl.Driver)
		err := nwd.Decode()
		if err != nil {
			return nil
		}

		intfName := nwd.Port.Name
		log.Debugf("init interface %s", intfName)
		reporters = append(reporters, func() ([]metrics, error) {
			nl := netlink.NewDecoder(intfName)
			log.Debugf("interface %s", intfName)
			err = nl.Decode()
			if err != nil {
				return nil, err
			}

			// TODO: hardware info like a product name
			labels := map[string]string{
				"interface": intfName,
			}

			points := map[string]interface{}{
				"nw_rx_packets": nl.Stats.RxPackets,
				"nw_tx_packets": nl.Stats.TxPackets,
				"nw_rx_bytes":   nl.Stats.RxBytes,
				"nw_tx_bytes":   nl.Stats.TxBytes,
				"nw_rx_errors":  nl.Stats.RxErrors,
				"nw_tx_errors":  nl.Stats.TxErrors,
				"nw_rx_dropped": nl.Stats.RxDropped,
				"nw_tx_dropped": nl.Stats.TxDropped,
			}

			return createMetricsList(labels, points)
		})
	}
	return reporters
}
