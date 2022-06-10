package pump

import (
	"time"

	"github.com/spf13/viper"
)

const (
	confKey = "pump"

	defaultPumpInterval = time.Duration(10 * time.Second)
	defaultOmitDetails  = true
)

// PumpOptions defines options for building a pump.
type PumpOptions struct {
	PumpInterval time.Duration `mapstructure:"pump-interval"`
	OmitDetails  bool          `mapstructure:"omit-details"`
}

func NewPumpOptions() *PumpOptions {
	opts := &PumpOptions{
		PumpInterval: defaultPumpInterval,
		OmitDetails:  defaultOmitDetails,
	}

	_ = viper.UnmarshalKey(confKey, opts)
	return opts
}
