package streamreceiver

import (
	"go.opentelemetry.io/collector/config"
)

// Config defines configuration for the receiver.
type Config struct {
	config.ReceiverSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct

	Organization string `mapstructure:"organization"`
	Project string `mapstructure:"project"`
	ApiKey string `mapstructure:"api_key"`
	StreamId string `mapstructure:"stream_id"`
	WindowSize string `mapstructure:"window_size"`
}