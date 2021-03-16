package main

import (
	"github.com/moxspec/moxspec/edac"
)

func scanMemory() []metricsReporter {
	return []metricsReporter{
		func() ([]metrics, error) {
			edacd := edac.NewDecoder()
			err := edacd.Decode()
			if err != nil {
				return nil, err
			}
			return makeMetricsEDACCtls(edacd.Controllers)
		},
	}
}

func makeMetricsEDACCtls(ctls []*edac.MemoryController) ([]metrics, error) {
	var mList []metrics
	for _, mc := range ctls {
		labels := map[string]string{
			"name": mc.Name,
		}

		points := map[string]interface{}{
			"edac_mc_size":            mc.Size,
			"edac_mc_ce_count":        mc.CECount,
			"edac_mc_ce_noinfo_count": mc.CENoInfoCount,
			"edac_mc_ue_count":        mc.UECount,
			"edac_mc_ue_noinfo_count": mc.UENoInfoCount,
		}

		ms, err := createMetricsList(labels, points)
		if err != nil {
			return nil, err
		}
		mList = append(mList, ms...)

		mList, err = appendMetricsEDACCSrows(mList, mc)
		if err != nil {
			return nil, err
		}
	}
	return mList, nil
}

func appendMetricsEDACCSrows(mList []metrics, mc *edac.MemoryController) ([]metrics, error) {
	for _, csrow := range mc.CSRows {
		labels := map[string]string{
			"name": csrow.Name,
		}

		points := map[string]interface{}{
			"edac_csrow_size":     csrow.Size,
			"edac_csrow_ce_count": csrow.CECount,
			"edac_csrow_ue_count": csrow.UECount,
		}

		ms, err := createMetricsList(labels, points)
		if err != nil {
			return nil, err
		}
		mList = append(mList, ms...)

		mList, err = appendMetricsEDACCSchans(mList, csrow)
		if err != nil {
			return nil, err
		}
	}
	return mList, nil
}

func appendMetricsEDACCSchans(mList []metrics, csrow *edac.ChipSelectRow) ([]metrics, error) {
	for _, c := range csrow.Channels {
		labels := map[string]string{
			"name":  c.Name,
			"label": c.Label,
		}

		points := map[string]interface{}{
			"edac_ch_ce_count": csrow.CECount,
		}

		ms, err := createMetricsList(labels, points)
		if err != nil {
			return nil, err
		}
		mList = append(mList, ms...)
	}
	return mList, nil
}
