package serviceexporter


import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)


const (
	// The value of "type" key in configuration.
	typeStr = "service"
)

// NewFactory creates a factory for OTLP exporter.
func NewFactory() component.ExporterFactory {
	return exporterhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		exporterhelper.WithTraces(createTracesExporter),
	)
}

func createDefaultConfig() config.Exporter {
	return &Config{
		ExporterSettings: config.NewExporterSettings(config.NewID(typeStr)),
	}
}

func createTracesExporter(
	_ context.Context,
	set component.ExporterCreateSettings,
	cfg config.Exporter,
) (component.TracesExporter, error) {
	oCfg := cfg.(*Config)
	se := NewServiceExporter(set.Logger, oCfg)

	return exporterhelper.NewTracesExporter(
		cfg,
		set,
		se.ConsumeTraces,
		exporterhelper.WithStart(se.Start),
		exporterhelper.WithShutdown(se.Shutdown),
	)
}