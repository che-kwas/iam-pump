package kafka

import (
	"github.com/spf13/viper"
)

const (
	confKey = "kafka"

	defaultVersion       = "3.2.0"
	defaultGroupID       = "iam-pump"
	defaultTopic         = "iam"
	defaultFromBeginning = true
	defaultPoolSize      = 10
)

// KafkaOptions defines options for building a kafka consumer.
type KafkaOptions struct {
	Brokers       []string
	Version       string
	GroupID       string `mapstructure:"group-id"`
	Topic         string
	FromBeginning bool `mapstructure:"from-beginning"`
	PoolSize      int  `mapstructure:"pool-size"`
}

func getKafkaOpts() (*KafkaOptions, error) {
	opts := &KafkaOptions{
		Version:       defaultVersion,
		GroupID:       defaultGroupID,
		Topic:         defaultTopic,
		FromBeginning: defaultFromBeginning,
		PoolSize:      defaultPoolSize,
	}

	if err := viper.UnmarshalKey(confKey, opts); err != nil {
		return nil, err
	}
	return opts, nil
}
