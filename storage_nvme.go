package main

import (
	"github.com/moxspec/moxspec/blk/nvme"
	"github.com/moxspec/moxspec/nvmeadm"
)

func scanNVMeController(path, driver string) ([]metricsReporter, error) {
	nvmed := nvme.NewDecoder(path)
	err := nvmed.Decode()
	if err != nil {
		return nil, err
	}
	deviceFile := "/dev/" + nvmed.Name

	return []metricsReporter{
		func() ([]metrics, error) {
			admd := nvmeadm.NewDecoder(deviceFile)
			err = admd.Decode()
			if err != nil {
				return nil, err
			}

			var mList []metrics

			labels := map[string]string{
				"driver":       driver,
				"name":         nvmed.Name,
				"path":         path,
				"serialNumber": admd.SerialNumber,
				"model":        admd.ModelNumber,
				"firmware":     admd.FirmwareRevision,
			}

			points := map[string]interface{}{
				"disk_cur_temp":              admd.CurTemp,
				"disk_warn_temp":             admd.WarnTemp,
				"disk_crit_temp":             admd.CritTemp,
				"disk_byte_read":             admd.ByteRead,
				"disk_byte_written":          admd.ByteWritten,
				"disk_size":                  admd.Size,
				"disk_power_cycle_count":     admd.PowerCycleCount,
				"disk_power_on_hours":        admd.PowerOnHours,
				"disk_unsafe_shutdown_count": admd.UnsafeShutdownCount,
			}

			ms, err := createMetricsList(labels, points)
			if err != nil {
				return nil, err
			}
			mList = append(mList, ms...)

			return appendMetricsNVMeNS(mList, admd, nvmed)
		},
	}, nil
}

func appendMetricsNVMeNS(mList []metrics, admd *nvmeadm.Device, nvmed *nvme.Controller) ([]metrics, error) {
	for _, n := range nvmed.Namespaces {
		sz := admd.GetNamespaceSize(n.ID())
		if sz <= 0 {
			sz = n.Size()
		}

		labels := map[string]string{
			"name":   n.Name,
			"parent": nvmed.Name,
		}

		points := map[string]interface{}{
			"drive_namespace_size":           sz,
			"drive_namespace_phy_block_size": n.PhyBlockSize,
			"drive_namespace_log_block_size": n.LogBlockSize,
		}

		ms, err := createMetricsList(labels, points)
		if err != nil {
			return nil, err
		}
		mList = append(mList, ms...)
	}
	return mList, nil
}
