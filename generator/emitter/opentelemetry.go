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
	"math/rand"
	"os"

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
	convertedSpans := make(map[[8]byte]oteltrace.Span)
	var traceId oteltrace.TraceID

	createOtelSpan := func(s *trace.Span) {
		oteltracer := e.getTracer(s.Service)

		var parentSpanId oteltrace.SpanID

		if len(s.Refs) > 0 {
			parentSpanId = s.Refs[0].FromSpanId
			traceId = convertedSpans[s.Refs[0].FromSpanId].SpanContext().TraceID()
		} else {
			var tid [16]byte
			rand.Read(tid[:])
			traceId = tid
		}

		parent := oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: traceId,
			SpanID: parentSpanId,
			Remote: true,
			TraceFlags: 0x1,
		})

		ctx := oteltrace.ContextWithRemoteSpanContext(context.Background(), parent)
		_, span := oteltracer.Start(ctx, s.OperationName)
		span.SetAttributes(semconv.HTTPMethodKey.String("GET"))
		span.SetAttributes(
			semconv.HTTPURLKey.String(fmt.Sprintf("http://%s%s", s.Service.ServiceName, s.OperationName)))

		for _, v := range s.Tags {
			span.SetAttributes(attribute.String(v.Key, v.Value))
		}

		convertedSpans[s.ID] = span
	}
	closeOtelSpan := func(s *trace.Span) {
		convertedSpans[s.ID].End()
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