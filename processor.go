package main

import (
	"fmt"

	"github.com/moxspec/moxspec-metrics-agent/rapl"
	"github.com/moxspec/moxspec/cpu"
	"github.com/moxspec/moxspec/msr"
)

func scanProcessor() []metricsReporter {
	return []metricsReporter{
		func() ([]metrics, error) {
			cput := cpu.NewDecoder()
			err := cput.Decode()
			if err != nil {
				return nil, err
			}

			return makeMetricsProc(cput.Packages())
		},
		func() ([]metrics, error) {
			if !rapl.IsRaplEnabled() {
				return nil, fmt.Errorf("rapl is disabled")
			}

			var mList []metrics
			for _, p := range rapl.ScanPackages() {
				labels := map[string]string{
					"package_name": p.Name,
					"path":         p.Path,
				}

				points := map[string]interface{}{
					"rapl_package_energy_uj": p.EnergyVal(),
				}

				ms, err := createMetricsList(labels, points)
				if err != nil {
					return nil, err
				}
				mList = append(mList, ms...)

				mList, err = appendMetricsRaplConsts(mList, p)
				if err != nil {
					return nil, err
				}
			}
			return mList, nil
		},
	}
}

func appendMetricsRaplConsts(mList []metrics, p rapl.Package) ([]metrics, error) {
	for _, c := range p.Consts {
		labels := map[string]string{
			"package_name": p.Name,
			"path":         p.Path,
			"const_name":   c.Name,
		}
		points := map[string]interface{}{
			"rapl_package_power_limit": c.PowerLimit,
			"rapl_package_max_power":   c.MaxPower,
			"rapl_package_time_window": c.TimeWindow,
		}
		ms, err := createMetricsList(labels, points)
		if err != nil {
			return nil, err
		}
		mList = append(mList, ms...)
	}
	return mList, nil
}

func makeMetricsProc(pkgs []*cpu.Package) ([]metrics, error) {
	var mList []metrics
	for _, p := range pkgs {
		for _, nd := range p.Nodes() {
			for _, cr := range nd.Cores() {
				msrd := msr.NewDecoder(cr.Threads[0], msr.INTEL)
				err := msrd.Decode()
				if err != nil {
					return nil, err
				}

				labels := map[string]string{
					"package": fmt.Sprintf("%d", p.ID),
					"node":    fmt.Sprintf("%d", nd.ID),
					"core":    fmt.Sprintf("%d", cr.ID),
				}

				points := map[string]interface{}{
					"cpu_throttlecount":    cr.ThrottleCount,
					"cpu_scaling_cur_freq": cr.Scaling.CurFreq,
					"cpu_scaling_max_freq": cr.Scaling.MaxFreq,
					"cpu_scaling_min_freq": cr.Scaling.MinFreq,
					"cpu_temp":             msrd.Temp,
				}

				ms, err := createMetricsList(labels, points)
				if err != nil {
					return nil, err
				}

				mList = append(mList, ms...)
			}
		}
	}
	return mList, nil
}
