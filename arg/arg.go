package arg

import (
	"flag"
)

var FlagConfig = struct {
	Levels []bool
	MinLvl int
}{
	Levels: make([]bool, 6),
	MinLvl: 0,
}

func init() {

	flag.BoolVar(&FlagConfig.Levels[0], "v", false, "Enable alert and error logging")
	flag.BoolVar(&FlagConfig.Levels[1], "vv", false, "Enable warn and lvl1 logging")
	flag.BoolVar(&FlagConfig.Levels[2], "vvv", false, "Enable highlight and lvl2 logging")
	flag.BoolVar(&FlagConfig.Levels[3], "vvvv", false, "Enable inform and lvl3 logging")
	flag.BoolVar(&FlagConfig.Levels[4], "vvvvv", false, "Enable log and lvl4 logging")
	flag.BoolVar(&FlagConfig.Levels[5], "vvvvvv", false, "Enable trace and lvl5 logging")

	flag.Parse()
	for i, c := range FlagConfig.Levels {
		if c {
			FlagConfig.MinLvl = i + 1
		}
	}
}
