package pump

import (
	"time"

	"github.com/spf13/viper"
)

const (
	confKey = "pump"

	defaultPumpInterval = time.Duration(10 * time.Second)
)

// PumpOptions defines options for building a pump.
type PumpOptions struct {
	PumpInterval time.Duration `mapstructure:"pump-interval"`
}

func NewPumpOptions() *PumpOptions {
	opts := &PumpOptions{
		PumpInterval: defaultPumpInterval,
	}

	_ = viper.UnmarshalKey(confKey, opts)
	return opts
}
