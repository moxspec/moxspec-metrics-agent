package rapl

import (
	"path/filepath"

	"github.com/moxspec/moxspec/loglet"
	"github.com/moxspec/moxspec/util"
)

var (
	log *loglet.Logger
)

func init() {
	log = loglet.NewLogger("rapl")
}

const (
	raplDir = "/sys/class/powercap/intel-rapl"
)

// IsRaplEnabled returns if rapl is enabled
func IsRaplEnabled() bool {
	stat, err := util.LoadByte(filepath.Join(raplDir, "enabled"))
	if err != nil {
		return false
	}

	if stat == 0 {
		return false
	}

	return true
}
