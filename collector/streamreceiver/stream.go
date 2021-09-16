package streamreceiver

import (
	"context"
	"fmt"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/model/pdata"
	"go.uber.org/zap"
	"net/url"
	"time"
)

type streamReceiver struct {
	logger         *zap.Logger
	traceConsumer  consumer.Traces
	ticker         *time.Ticker
	tickerDuration time.Duration
	client         Client
	stop           chan struct{}
}

func (s streamReceiver) Start(ctx context.Context, host component.Host) error {
	go s.consumeStreamData()

	s.ticker = time.NewTicker(s.tickerDuration)
	s.stop = make(chan struct{})
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.consumeStreamData()
			case <-s.stop:
				s.ticker.Stop()
				return
			}
		}
	}()

	return nil
}

func convertStringToTraceId(traceId string) pdata.TraceID {
	var newTraceId [16]byte
	copy(newTraceId[:], traceId)
	return pdata.NewTraceID(newTraceId)
}

func convertStringToSpanId(spanId string) pdata.SpanID {
	var newSpanId [8]byte
	copy(newSpanId[:], spanId)
	return pdata.NewSpanID(newSpanId)
}

func convertTimeFromString(t int64) time.Time {
	return time.Unix(0, t*int64(time.Microsecond))
}

func (s streamReceiver) convertTrace(trace LightstepTrace) *pdata.Traces {
	traces := pdata.NewTraces()
	var reporterLookup = make(map[string]map[string]interface{})
	for _, reporter := range trace.Relationships.Reporters {
		reporterLookup[reporter.ReporterID] = reporter.Attributes
	}

	rspanSlice := traces.ResourceSpans()

	for _, span := range trace.Attributes.Spans {
		// add resource metadata
		rspan := rspanSlice.AppendEmpty()
		resource := rspan.Resource()
		resourceAttrs := reporterLookup[span.ReporterID]
		for k, v := range resourceAttrs {
			resource.Attributes().InsertString(k, fmt.Sprintf("%s", v))
		}
		// add span attrs
		ils := rspan.InstrumentationLibrarySpans().AppendEmpty()
		spans := ils.Spans()

		otelSpan := spans.AppendEmpty()
		otelSpan.SetTraceID(convertStringToTraceId(span.TraceID))
		otelSpan.SetSpanID(convertStringToSpanId(span.SpanID))

		otelSpan.SetStartTimestamp(pdata.NewTimestampFromTime(convertTimeFromString(span.StartTimeMicros)))
		otelSpan.SetEndTimestamp(pdata.NewTimestampFromTime(convertTimeFromString(span.EndTimeMicros)))
		otelSpan.SetKind(pdata.SpanKindUnspecified)
		otelSpan.SetName(span.SpanName)
		for k, v := range span.Tags {
			if k == "span.kind" && v == "server" {
				otelSpan.SetKind(pdata.SpanKindServer)
			} else if k == "span.kind" && v == "client" {
				otelSpan.SetKind(pdata.SpanKindClient)
			}

			if k == "parent_span_guid" {
				otelSpan.SetParentSpanID(convertStringToSpanId(fmt.Sprintf("%s", v)))
			}
			otelSpan.Attributes().InsertString(k, fmt.Sprintf("%s", v))
		}
	}
	return &traces
}

func (s streamReceiver) consumeStreamData() error {
	traces, err := s.getTraces()
	if err != nil {
		s.logger.Error("Could not get traces", zap.Error(err))
	}

	for _, t := range traces {
		otelTrace := s.convertTrace(t)
		_ = s.traceConsumer.ConsumeTraces(context.Background(), *otelTrace)
	}

	return nil
}

func (s streamReceiver) getTraces() ([]LightstepTrace, error) {
	var traces []LightstepTrace
	resp, err := s.client.GetStreamTraces()
	if err != nil {
		s.logger.Info(fmt.Sprintf("Could not get traces: %v", err))
		return nil, err
	}

	exemplars := resp.Data.Attributes.Exemplars
	s.logger.Info("Found trace exemplars",
		zap.String("stream_id", s.client.StreamID()),
		zap.Int("exemplar_count", len(exemplars)))

	for _, exemplar := range exemplars {
		traceResp, err := s.client.GetTrace(exemplar.SpanGUID)
		if err != nil {
			s.logger.Info(fmt.Sprintf("Could not get trace: %v", err))
			continue
		}
		t := traceResp.Data[0]
		traces = append(traces, t)
	}

	return traces, nil
}

func (s streamReceiver) Shutdown(ctx context.Context) error {
	close(s.stop)
	return nil
}


func newTraceReceiver(config *Config,
	consumer consumer.Traces,
	logger *zap.Logger) (component.TracesReceiver, error) {
	var sReceiver = streamReceiver{}

	if consumer == nil {
		return nil, componenterror.ErrNilNextConsumer
	}
	u, _ := url.Parse("https://api.lightstep.com/public/v0.2/")
	sReceiver.logger = logger
	c := NewClientProvider(*u, config.Organization, config.Project, config.ApiKey, config.WindowSize, config.StreamId, logger).BuildClient()
	sReceiver.client = c
	sReceiver.traceConsumer = consumer

	var tickerDuration time.Duration
	tickerDuration, err := time.ParseDuration(config.WindowSize)
	if err != nil {
		tickerDuration = 5 * time.Minute
	}
	sReceiver.tickerDuration = tickerDuration

	return &sReceiver, nil
}
