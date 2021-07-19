package emitter

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/smithclay/synthetic-load-generator-go/topology"
	"github.com/smithclay/synthetic-load-generator-go/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"

	"os"
	"time"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"

	"log"

	metricexport "go.opentelemetry.io/otel/sdk/export/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	ValueRecorderType = "ValueRecorder"
	CounterType       = "Counter"
)

type OpenTelemetryEmitter struct {
	collectorUrl                string
	flushIntervalMillis         int
	stdout                      bool
	serviceNameToTracerProvider map[string]*sdktrace.TracerProvider
	serviceNameToMeterProvider  map[string]*metric.MeterProvider
}

func NewOpenTelemetryGrpcEmitter(collectorUrl string) *OpenTelemetryEmitter {
	return &OpenTelemetryEmitter{
		serviceNameToTracerProvider: make(map[string]*sdktrace.TracerProvider),
		serviceNameToMeterProvider:  make(map[string]*metric.MeterProvider),
		collectorUrl:                collectorUrl,
		stdout:                      false,
	}
}

func NewOpenTelemetryStdoutEmitter() *OpenTelemetryEmitter {
	return &OpenTelemetryEmitter{
		serviceNameToTracerProvider: make(map[string]*sdktrace.TracerProvider),
		serviceNameToMeterProvider:  make(map[string]*metric.MeterProvider),
		stdout:                      true,
	}
}

func (e *OpenTelemetryEmitter) EmitMetric(metrics []topology.Metric, service string) {
	meter := e.getMeter(service)
	for _, m := range metrics {
		if m.Type == ValueRecorderType {
			recorder, err := meter.NewFloat64ValueRecorder(m.Name,
				metric.WithDescription("Synthetic metric via Lightstep Partner Toolkit"))
			if err != nil {
				log.Fatalf("error creating recorder: %v", err)
			}
			recorder.Record(context.Background(), m.Min+rand.Float64()*(m.Max-m.Min),
				attribute.String("service.name", service),
				attribute.Bool("synthetic", true))
		} else if m.Type == CounterType {
			counter, err := meter.NewFloat64Counter(m.Name,
				metric.WithDescription("Synthetic metric via Lightstep Partner Toolkit"))
			if err != nil {
				log.Fatalf("error creating counter: %v", err)
			}
			counter.Add(context.Background(), m.Min+rand.Float64()*(m.Max-m.Min),
				attribute.String("service.name", service),
				attribute.Bool("synthetic", true))
		}
	}
}

func (e *OpenTelemetryEmitter) EmitTrace(t *trace.Trace) {
	convertedSpans := make(map[oteltrace.SpanID]oteltrace.Span)
	spanContext := make(map[oteltrace.SpanID]context.Context)

	createOtelSpan := func(s *trace.Span) {
		t := e.getTracer(s.Service)

		var parentCtx context.Context
		if len(s.Refs) == 0 {
			parentCtx = context.TODO()
		} else {
			parentCtx = oteltrace.ContextWithRemoteSpanContext(spanContext[s.Refs[0].FromSpanId],
				convertedSpans[s.Refs[0].FromSpanId].SpanContext())
		}

		ctx, span := t.Start(parentCtx, s.OperationName, oteltrace.WithTimestamp(time.Unix(0, s.StartTimeMicros)))
		span.SetAttributes(semconv.HTTPMethodKey.String("GET"))
		span.SetAttributes(
			semconv.HTTPURLKey.String(fmt.Sprintf("http://%s%s", s.Service.ServiceName, s.OperationName)))

		for _, v := range s.Tags {
			span.SetAttributes(attribute.String(v.Key, v.Value))
		}
		convertedSpans[s.ID] = span
		spanContext[s.ID] = ctx
	}
	closeOtelSpan := func(s *trace.Span) {
		convertedSpans[s.ID].End(oteltrace.WithTimestamp(time.Unix(0, s.EndTimeMicros)))
	}
	prePostOrder(t, createOtelSpan, closeOtelSpan)
}

func (e *OpenTelemetryEmitter) Close() {
	for _, tp := range e.serviceNameToTracerProvider {
		_ = tp.ForceFlush(context.Background())
		_ = tp.Shutdown(context.Background())
	}
}

func (e *OpenTelemetryEmitter) getTracer(service trace.Service) oteltrace.Tracer {
	if _, ok := e.serviceNameToTracerProvider[service.ServiceName]; !ok {
		e.serviceNameToTracerProvider[service.ServiceName] = initTracer(service.ServiceName, e.stdout, e.collectorUrl)
	}
	tp := e.serviceNameToTracerProvider[service.ServiceName]

	return tp.Tracer(service.ServiceName)
}

func (e *OpenTelemetryEmitter) getMeter(service string) metric.Meter {
	if _, ok := e.serviceNameToMeterProvider[service]; !ok {
		e.serviceNameToMeterProvider[service] = initMeter(service, e.stdout, e.collectorUrl)
	}
	mp := e.serviceNameToMeterProvider[service]

	return (*mp).Meter(service)
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

const LightstepPublicIngest = "ingest.lightstep.com:443"

func initMeter(serviceName string, isStdout bool, collectorUrl string) *metric.MeterProvider {
	var exp metricexport.Exporter
	var err error

	if isStdout {
		exp, err = stdoutmetric.New(stdoutmetric.WithPrettyPrint())
		if err != nil {
			log.Fatalf("creating stdoutmetric exporter: %v", err)
		}
	} else {
		client := otlpmetricgrpc.NewClient(
			otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
			otlpmetricgrpc.WithEndpoint(collectorUrl),
			otlpmetricgrpc.WithHeaders(map[string]string{
				"lightstep-access-token": os.Getenv("LS_ACCESS_TOKEN"),
			}),
		)
		exp, err = otlpmetric.New(context.Background(), client)
		if err != nil {
			log.Fatalf("creating otlpmetricgrpc exporter: %v", err)
		}
	}

	pusher := controller.New(
		processor.New(
			simple.NewWithInexpensiveDistribution(),
			exp,
		),
		controller.WithExporter(exp),
		controller.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
			),
		))
	if err := pusher.Start(context.Background()); err != nil {
		log.Fatalf("starting push controller: %v", err)
	}
	mp := pusher.MeterProvider()
	return &mp
}

func initTracer(serviceName string, isStdout bool, collectorUrl string) *sdktrace.TracerProvider {
	var err error
	var exp sdktrace.SpanExporter

	if isStdout {
		exp, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			log.Panicf("failed to initialize stdout exporter %v\n", err)
			return nil
		}
	} else {
		ctx := context.Background()
		exp, err = otlptrace.New(
			ctx,
			otlptracegrpc.NewClient(
				otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
				otlptracegrpc.WithEndpoint(collectorUrl),
				otlptracegrpc.WithTimeout(10*time.Second),
				otlptracegrpc.WithHeaders(map[string]string{
					"lightstep-access-token": os.Getenv("LS_ACCESS_TOKEN"),
				}),
			),
		)
		if err != nil {
			log.Panicf("count not initialize otlptrace exporter: %v", err)
			return nil
		}
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
			),
		))
	return tp
}
