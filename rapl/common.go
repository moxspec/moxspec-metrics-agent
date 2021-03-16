package rapl

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/moxspec/moxspec/util"
)

// Package represents a rapl domain
type Package struct {
	Domain
	Dram Domain
}

// Domain represents a rapl domain
type Domain struct {
	Name   string
	Path   string
	Consts []Const
}

// EnergyFile returns a energy file path
func (d Domain) EnergyFile() string {
	return filepath.Join(d.Path, "energy_uj")
}

// EnergyVal returns the current energy value
func (d Domain) EnergyVal() float64 {
	val, _ := util.LoadUint64(d.EnergyFile())
	return float64(val)
}

// Const represents a rapl constraint
type Const struct {
	Prefix     string `json:"prefix"`
	Name       string `json:"name"`
	MaxPower   uint64 `json:"maxPower"`
	PowerLimit uint64 `json:"powerLimit"`
	TimeWindow uint64 `json:"timeWindow"`
}

// PowerLimitName returns the powerlimit const name
func (c Const) PowerLimitName() string {
	return c.Name + "_power_limit_uw"
}

// MaxPowerName returns the maxpower const name
func (c Const) MaxPowerName() string {
	return c.Name + "_max_power_uw"
}

// TimeWindowName returns the timewindow const name
func (c Const) TimeWindowName() string {
	return c.Name + "_time_window_us"
}

// ScanPackages returns rapl package info by scanning from the given path
func ScanPackages() []Package {
	var pkgs []Package
	for _, f := range util.FilterPrefixedDirs(raplDir, "intel-rapl:") {
		pkgs = append(pkgs, scanDomains(f))
	}
	return pkgs
}

func scanDomains(path string) Package {
	name, _ := util.LoadString(filepath.Join(path, "name"))
	consts := scanConsts(path)

	dom := Domain{
		Name:   name,
		Path:   path,
		Consts: consts,
	}

	dramPath := filepath.Join(path, fmt.Sprintf("%s:%s", filepath.Base(path), "0"))
	dramName, _ := util.LoadString(filepath.Join(dramPath, "name"))
	dramConsts := scanConsts(dramPath)

	dram := Domain{
		Name:   dramName,
		Path:   dramPath,
		Consts: dramConsts,
	}

	return Package{
		dom,
		dram,
	}
}

func scanConsts(path string) []Const {
	dict := make(map[string]*Const)
	for _, f := range util.FilterPrefixedFiles(path, "constraint_") {
		// e.g:
		//   constraint_0_name
		//   constraint_0_time_window_us
		fname := filepath.Base(f)
		elm := strings.Split(fname, "_")

		if len(elm) < 3 {
			continue // something wrong
		}

		// e.g:
		//   constraint_0_time_window_us
		prefix := strings.Join(elm[:2], "_") // constraint_0
		key := strings.Join(elm[2:], "_")    // time_window_us

		if _, ok := dict[prefix]; !ok {
			r := new(Const)
			r.Prefix = prefix
			dict[prefix] = r
		}

		switch key {
		case "name":
			dict[prefix].Name, _ = util.LoadString(f)
		case "max_power_uw":
			dict[prefix].MaxPower, _ = util.LoadUint64(f)
		case "power_limit_uw":
			dict[prefix].PowerLimit, _ = util.LoadUint64(f)
		case "time_window_us":
			dict[prefix].TimeWindow, _ = util.LoadUint64(f)
		}
	}

	var consts []Const
	for _, v := range dict {
		consts = append(consts, *v)
	}
	return consts
}
