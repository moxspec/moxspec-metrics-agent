package main

import (
	"fmt"

	"github.com/moxspec/moxspec/blk/raid"
	"github.com/moxspec/moxspec/pci"
	"github.com/moxspec/moxspec/raidcli/megacli"
	"github.com/moxspec/moxspec/spc"
	"github.com/moxspec/moxspec/spc/megaraid"
)

func scanMegaRAIDController(ctl pci.Device, logDrives []*raid.VirtualDisk) ([]metricsReporter, error) {
	if !megacli.Available() {
		return nil, fmt.Errorf("megacli is not available")
	}

	mctls, err := megacli.GetControllers()
	if err != nil {
		return nil, fmt.Errorf("do not parse megaraid: %w", err)
	}

	var mctl *megacli.Controller
	for _, m := range mctls {
		if m.Bus == ctl.Bus && m.Device == ctl.Device && m.Function == ctl.Function {
			mctl = m
			break
		}
	}

	if mctl == nil {
		return nil, fmt.Errorf("the controller not found in megacli")
	}

	err = mctl.Decode()
	if err != nil {
		return nil, fmt.Errorf("failed to decode megacli: %w", err)
	}

	// label base
	baseLabels := map[string]string{
		"product_name":  mctl.ProductName,
		"bios":          mctl.BIOS,
		"firmware":      mctl.Firmware,
		"serial_number": mctl.SerialNumber,
		"adapter_id":    fmt.Sprintf("%d", mctl.Number),
		"battery":       fmt.Sprintf("%t", mctl.Battery),
	}

	var mReporters []metricsReporter

	// ctl.LogDrives may contain pass-through drives that must be separated to ctl.PassthroughDrives.
	for _, ldrv := range logDrives {
		log.Debugf("scanning megacli data for %s", ldrv.Path)

		log.Debugf("wwn: %s", ldrv.WWN)
		if ldrv.WWN != "" {
			log.Debug("this logical drive has wwn. possibly it is pass-through disk (jbod)")
			log.Debug("scanning wwn")

			ptpd := mctl.GetPTPhyDriveByWWN(ldrv.WWN)
			if ptpd != nil {
				log.Debugf("[enc:slt] = %s:%s", ptpd.EnclosureID, ptpd.SlotNumber)
				mReporters = append(mReporters, makeMegaRAIDpdReporter(ptpd, mctl.Number, baseLabels))
				continue
			}

			log.Debug("cound not find wwn from controller")
			log.Debug("continue scanning")
		}

		ld := mctl.GetLogDriveByTarget(ldrv.Target)
		if ld == nil {
			continue
		}

		log.Debugf("megacli has the ld data for the target %d", ldrv.Target)
		ldrvBaseLabels := copyLabels(baseLabels)
		ldrvBaseLabels["raid_lv"] = string(ld.RAIDLv)
		ldrvBaseLabels["cache_policy"] = ld.CachePolicy
		ldrvBaseLabels["status"] = ld.State
		ldrvBaseLabels["stripe_size"] = fmt.Sprintf("%d", ld.StripSize)
		ldrvBaseLabels["group_label"] = ld.Label

		for _, pd := range ld.PhyDrives {
			log.Debugf("found phy drive: %s", pd.Model)
			mReporters = append(mReporters, makeMegaRAIDpdReporter(pd, mctl.Number, ldrvBaseLabels))
		}
	}

	for _, pd := range mctl.UnconfDrives {
		log.Debugf("found unconfigured phy drive: %s", pd.Model)
		mReporters = append(mReporters, makeMegaRAIDpdReporter(pd, mctl.Number, baseLabels))
	}

	return mReporters, nil
}

func makeMegaRAIDpdReporter(phyDrv *megacli.PhyDrive, ctlNum int, baseLabels map[string]string) metricsReporter {
	pd := *phyDrv
	labels := copyLabels(baseLabels)

	return func() ([]metrics, error) {
		d := megaraid.NewDecoder(ctlNum, int(pd.DeviceID), spc.CastDiskType(pd.Type))
		err := d.Decode()
		if err != nil {
			return nil, err
		}

		blockSize := uint64(pd.LogBlockSize)
		if pd.Type == "SAS" {
			blockSize = 1
		}

		labels["enclosure"] = pd.EnclosureID
		labels["slot"] = pd.SlotNumber
		labels["status"] = pd.State
		labels["firmware"] = d.FirmwareRevision
		labels["serial_number"] = d.SerialNumber
		labels["model"] = d.ModelNumber
		labels["neg_speed"] = pd.DriveSpeed
		labels["transport"] = pd.Type

		points := map[string]interface{}{
			"disk_size":                  pd.Size,
			"disk_error_count":           pd.MediaErrorCount,
			"disk_byte_written":          uint64(d.TotalLBAWritten) * blockSize,
			"disk_byte_read":             uint64(d.TotalLBARead) * blockSize,
			"disk_power_cycle_count":     uint64(d.TotalLBARead) * blockSize,
			"disk_power_on_hours":        d.PowerOnHours,
			"disk_unsafe_shutdown_count": d.UnsafeShutdownCount,
		}

		return createMetricsList(labels, points)

		// TODO: smart
		/*
			for _, rec := range d.ErrorRecords {
				e := new(model.SMARTRecord)
				e.ID = rec.ID
				e.Current = rec.Current
				e.Worst = rec.Worst
				e.Raw = rec.Raw
				e.Threshold = rec.Threshold
				e.Name = rec.Name
				p.ErrorRecords = append(p.ErrorRecords, e)
			}
		*/
	}
}
