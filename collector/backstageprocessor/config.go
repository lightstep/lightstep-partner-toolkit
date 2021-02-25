package backstageprocessor

import (
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configmodels"
)

// Config defines configuration for http forwarder extension.
type Config struct {
	configmodels.ProcessorSettings `mapstructure:",squash"`

	BackstageServer confighttp.HTTPClientSettings `mapstructure:"backstage_server"`
}
