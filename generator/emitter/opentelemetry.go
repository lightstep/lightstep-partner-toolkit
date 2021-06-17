package emitter

import (
	"context"
	"fmt"
	"github.com/smithclay/synthetic-load-generator-go/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"
	"os"
	"time"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
)

type OpenTelemetryEmitter struct {
	collectorUrl string
	flushIntervalMillis int
	stdout bool
	serviceNameToTracerProvider map[string]*sdktrace.TracerProvider
}

func NewOpenTelemetryGrpcEmitter(collectorUrl string) *OpenTelemetryEmitter {
	return &OpenTelemetryEmitter{
		serviceNameToTracerProvider: make(map[string]*sdktrace.TracerProvider),
		collectorUrl: collectorUrl,
		stdout: false,
	}
}

func NewOpenTelemetryStdoutEmitter() *OpenTelemetryEmitter {
	return &OpenTelemetryEmitter{
		serviceNameToTracerProvider: make(map[string]*sdktrace.TracerProvider),
		stdout: true,
	}
}

func (e *OpenTelemetryEmitter) Emit(t *trace.Trace) {
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
		e.serviceNameToTracerProvider[service.ServiceName] = initTracer(service.ServiceName, e.stdout)
	}
	tp := e.serviceNameToTracerProvider[service.ServiceName]

	return tp.Tracer(service.ServiceName)
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func initTracer(serviceName string, isStdout bool) *sdktrace.TracerProvider {
	var err error
	var exp sdktrace.SpanExporter

	if isStdout {
		exp, err = stdout.NewExporter(stdout.WithPrettyPrint())
		if err != nil {
			log.Panicf("failed to initialize stdout exporter %v\n", err)
			return nil
		}
	} else {
		ctx := context.Background()
		exp, err = otlp.NewExporter(
			ctx,
			otlpgrpc.NewDriver(
				otlpgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
				otlpgrpc.WithEndpoint(getenv("OTEL_EXPORTER_OTLP_SPAN_ENDPOINT", "ingest.lightstep.com:443")),
				otlpgrpc.WithHeaders(map[string]string{
					"lightstep-access-token":    os.Getenv("LS_ACCESS_TOKEN"),
				}),
			),
		)
	}

	if err != nil {
		log.Fatalf("failed to initialize otelgrpc pipeline: %v", err)
	}
	bsp := sdktrace.NewBatchSpanProcessor(exp)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithResource(
			resource.NewWithAttributes(semconv.ServiceNameKey.String(serviceName)),
	))
	return tp
}