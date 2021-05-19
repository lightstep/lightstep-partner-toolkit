package backstageprocessor

import (
	"context"
	"fmt"
	"go.opentelemetry.io/collector/config"
	"net/url"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor/processorhelper"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/consumer/pdata"
)

const (
	// The value of "type" key in configuration.
	typeStr = "backstage"
)

var processorCapabilities = consumer.Capabilities{MutatesData: true}

// TODO: support one attrProc per service
type resourceProcessor struct {
	attrProcs map[string]*processorhelper.AttrProc
}

// ProcessTraces implements the TProcessor interface
func (rp *resourceProcessor) ProcessTraces(_ context.Context, td pdata.Traces) (pdata.Traces, error) {
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		resource := rss.At(i).Resource()
		attrs := resource.Attributes()

		serviceAttr, ok := attrs.Get("service.name")
		if !ok {
			continue
		}

		attrProc, ok := rp.attrProcs[serviceAttr.StringVal()]
		if !ok {
			continue
		}
		fmt.Printf("Adding service catalog metadata for service: %v\n", serviceAttr)
		attrProc.Process(attrs)
	}
	return td, nil
}

// NewFactory returns a new factory for the Resource processor.
func NewFactory() component.ProcessorFactory {
	return processorhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		processorhelper.WithTraces(createTraceProcessor),
	)
}

// Note: This isn't a valid configuration because the processor would do no work.
func createDefaultConfig() config.Processor {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(config.NewID(typeStr)),
	}
}

func createTraceProcessor(
	_ context.Context,
	params component.ProcessorCreateParams,
	cfg config.Processor,
	nextConsumer consumer.Traces) (component.TracesProcessor, error) {
	attrProcs, err := createAttrProcessor(cfg.(*Config), params.Logger)
	if err != nil {
		return nil, err
	}
	return processorhelper.NewTracesProcessor(
		cfg,
		nextConsumer,
		&resourceProcessor{attrProcs: attrProcs},
		processorhelper.WithCapabilities(processorCapabilities))
}

func createAttrProcessor(cfg *Config, logger *zap.Logger) (map[string]*processorhelper.AttrProc, error) {
	attrProcs := make(map[string]*processorhelper.AttrProc)

	// TODO: validate backstage service catalog api endpoint
	if len(cfg.BackstageServer.Endpoint) == 0 {
		return nil, fmt.Errorf("error creating \"%q\" processor due to missing required field \"backstage_endpoint\"", cfg.ProcessorSettings.ID())
	}
	epURL, err := url.Parse(cfg.BackstageServer.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("error creating \"%q\" processor due to bad server url in \"backstage_endpoint\": %v", cfg.ProcessorSettings.ID(), err)
	}

	cp := NewClientProvider(*epURL, logger)
	backstageClient := cp.BuildClient()
	entities, err := backstageClient.GetEntities()
	if err != nil {
		return nil, fmt.Errorf("error \"%q\" processor due to bad backstage server response: %v", cfg.ProcessorSettings.ID(), err)
	}

	for _, entity := range *entities {
		serviceName, ok := entity.Metadata.Annotations["opentelemetry.io/service-name"]
		if !ok {
			continue
		}

		fmt.Printf("Found service: %v\n", serviceName)

		attrProc, err := processorhelper.NewAttrProc(&processorhelper.Settings{
			Actions: []processorhelper.ActionKeyValue{
				{
					Key:    "service.catalog.uuid",
					Value:  entity.Metadata.UID,
					Action: processorhelper.UPSERT,
				},
			},
		})

		if err != nil {
			return nil, fmt.Errorf("error creating \"%q\" processor: %w", cfg.ProcessorSettings.ID(), err)
		}
		attrProcs[fmt.Sprintf("%v", serviceName)] = attrProc
	}

	return attrProcs, nil
}
