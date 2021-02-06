package shell

import (
	"strings"

	flag "github.com/spf13/pflag"
)

var arguments struct {
	deviceType  *string
	deviceIndex *int
	frequency   *string
}

func init() {
	arguments.deviceType = flag.StringP("device-type", "t", "", "Device type (rtlsdr, limesdr, hackrf)")
	arguments.deviceIndex = flag.IntP("device-index", "i", 0, "Device index")
	arguments.frequency = flag.StringP("frequency", "f", "1090", "Frequency (MHz)")

	flag.Parse()

	if *arguments.deviceIndex < 0 {
		*arguments.deviceIndex = 0
	}

	*arguments.deviceType = strings.ToLower(*arguments.deviceType)
}
