package main

import (
	"fmt"

	"github.com/moxspec/moxspec/blk/ahci"
	"github.com/moxspec/moxspec/spc"
	"github.com/moxspec/moxspec/spc/acs"
)

func scanAHCIController(path, driver string) ([]metricsReporter, error) {
	ahcid := ahci.NewDecoder(path)
	err := ahcid.Decode()
	if err != nil {
		return nil, err
	}

	var reporters []metricsReporter
	for _, disk := range ahcid.Disks {
		reporters = append(reporters, func() ([]metrics, error) {
			var mList []metrics

			acsd := acs.NewDecoder("/dev/" + disk.Name)
			err := acsd.Decode()
			if err != nil {
				return nil, err
			}

			labels := map[string]string{
				"path":          disk.Path,
				"driver":        disk.Driver,
				"device_name":   disk.Name,
				"product_model": acsd.ModelNumber,
				"form_factor":   acsd.FormFactor,
				"firmware":      acsd.FirmwareRevision,
				"serial_number": acsd.SerialNumber,
				"transport":     acsd.Transport,
				//"disk_self_test":             acsd.SelfTestSupport,
				//"disk_error_logging":         acsd.ErrorLoggingSupport,
			}

			points := map[string]interface{}{
				"disk_cur_temp":              acsd.CurTemp,
				"disk_max_temp":              acsd.MaxTemp,
				"disk_min_temp":              acsd.MinTemp,
				"disk_rotation":              acsd.Rotation,
				"disk_byte_read":             uint64(acsd.TotalLBARead) * uint64(disk.LogBlockSize),
				"disk_byte_written":          uint64(acsd.TotalLBAWritten) * uint64(disk.LogBlockSize),
				"disk_neg_speed":             acsd.NegSpeed,
				"disk_sig_speed":             acsd.SigSpeed,
				"disk_power_cycle_count":     acsd.PowerCycleCount,
				"disk_power_on_hours":        acsd.PowerOnHours,
				"disk_unsafe_shutdown_count": acsd.UnsafeShutdownCount,
			}
			ms, err := createMetricsList(labels, points)
			if err != nil {
				return nil, err
			}
			mList = append(mList, ms...)

			return appendMetricsAHCIErrRecords(mList, acsd)
		})
	}
	return reporters, nil
}

func appendMetricsAHCIErrRecords(mList []metrics, acsd *spc.Device) ([]metrics, error) {
	for _, rec := range acsd.ErrorRecords {
		labels := map[string]string{
			"product_model":  acsd.ModelNumber,
			"firmware":       acsd.FirmwareRevision,
			"serial_number":  acsd.SerialNumber,
			"attribute_id":   fmt.Sprintf("%d", rec.ID),
			"attribute_name": rec.Name,
		}
		points := map[string]interface{}{
			"disk_smart_current":   rec.Current,
			"disk_smart_worst":     rec.Worst,
			"disk_smart_raw":       rec.Raw,
			"disk_smart_threshold": rec.Threshold,
		}
		ms, err := createMetricsList(labels, points)
		if err != nil {
			return nil, err
		}
		mList = append(mList, ms...)
	}
	return mList, nil
}
