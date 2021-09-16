package internal

import (
	"fmt"
	"github.com/google/pprof/profile"
)

// CPUUsage takes in a byte array of a GZip-encoded .proto file of a golang runtime CPU Profile
// and returns the total CPU usage as a float: 1.0 is 100% of 1 core being used
func CPUUsage(b []byte) (float64, error) {
	prof, err := profile.ParseData(b)
	if err != nil {
		return 0, err
	}

	if prof.PeriodType.Type != "cpu" {
		return 0, NotCPUProfile{s: prof.PeriodType.Type}
	}

	var total int64
	for _, sample := range prof.Sample {
		total += sample.Value[1]
	}

	return float64(total) / float64(prof.DurationNanos), nil
}

type NotCPUProfile struct {
	s string
}

func (n NotCPUProfile) Error() string {
	return fmt.Sprintf("Profile type %s is not expected type of cpu", n.s)
}