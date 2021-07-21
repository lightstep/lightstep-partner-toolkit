package serviceexporter

import (
	"errors"
	"go.opentelemetry.io/collector/config/confighttp"

	"go.opentelemetry.io/collector/config"
)

// Config defines configuration for file exporter.
type Config struct {
	config.ExporterSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct

	Scraper confighttp.HTTPServerSettings `mapstructure:"scraper"`
}

var _ config.Exporter = (*Config)(nil)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {
	if cfg.Scraper.Endpoint == "" {
		return errors.New("endpoint must be non-empty")
	}

	return nil
}