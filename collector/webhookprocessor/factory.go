// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package webhookprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

type Option func(forwarder *httpServer) error

var processorCapabilities = component.ProcessorCapabilities{MutatesConsumedData: true}

const (
	// The value of extension "type" in configuration.
	typeStr configmodels.Type = "oteltest"

	// Default endpoints to bind to.
	defaultTracesEndpoint  = ":7070"
	defaultMetricsEndpoint = ":7071"
)

// NewFactory creates a factory for HostObserver extension.
func NewFactory() component.ProcessorFactory {
	return processorhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		processorhelper.WithTraces(createTraceProcessor),
		processorhelper.WithMetrics(createMetricsProcessor))
}

func createDefaultConfig() configmodels.Processor {
	return &Config{
		ProcessorSettings: configmodels.ProcessorSettings{
			TypeVal: typeStr,
			NameVal: string(typeStr),
		},
		TracesIngress: confighttp.HTTPServerSettings{
			Endpoint: defaultTracesEndpoint,
		},
		MetricsIngress: confighttp.HTTPServerSettings{
			Endpoint: defaultMetricsEndpoint,
		},
	}
}

func createProcessor(
	params component.ProcessorCreateParams,
	cfg configmodels.Processor,
	serverType string,
	options ...Option,
) (*httpServer, error) {
	oCfg := cfg.(*Config)
	server, err := newHTTPServer(oCfg, params.Logger, serverType)
	return server, err
}

func createTraceProcessorWithOptions(
	_ context.Context,
	params component.ProcessorCreateParams,
	cfg configmodels.Processor,
	next consumer.TracesConsumer,
	options ...Option,
) (component.TracesProcessor, error) {
	kp, err := createProcessor(params, cfg, TracesServer, options...)
	if err != nil {
		return nil, err
	}

	return processorhelper.NewTraceProcessor(
		cfg,
		next,
		kp,
		processorhelper.WithCapabilities(processorCapabilities),
		processorhelper.WithStart(kp.Start),
		processorhelper.WithShutdown(kp.Shutdown))
}

func createTraceProcessor(ctx context.Context, params component.ProcessorCreateParams, cfg configmodels.Processor, next consumer.TracesConsumer) (component.TracesProcessor, error) {
	return createTraceProcessorWithOptions(ctx, params, cfg, next)
}

func createMetricsProcessorWithOptions(
	_ context.Context,
	params component.ProcessorCreateParams,
	cfg configmodels.Processor,
	nextMetricsConsumer consumer.MetricsConsumer,
	options ...Option,
) (component.MetricsProcessor, error) {
	kp, err := createProcessor(params, cfg, MetricsServer, options...)
	if err != nil {
		return nil, err
	}

	return processorhelper.NewMetricsProcessor(
		cfg,
		nextMetricsConsumer,
		kp,
		processorhelper.WithCapabilities(processorCapabilities),
		processorhelper.WithStart(kp.Start),
		processorhelper.WithShutdown(kp.Shutdown))
}
func createMetricsProcessor(
	ctx context.Context,
	params component.ProcessorCreateParams,
	cfg configmodels.Processor,
	nextMetricsConsumer consumer.MetricsConsumer,
) (component.MetricsProcessor, error) {
	return createMetricsProcessorWithOptions(ctx, params, cfg, nextMetricsConsumer)
}
