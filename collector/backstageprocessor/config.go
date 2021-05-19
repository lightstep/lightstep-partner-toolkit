package backstageprocessor

import (
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confighttp"
)

// Config defines configuration for http forwarder extension.
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"`

	BackstageServer confighttp.HTTPClientSettings `mapstructure:"backstage_server"`
}
