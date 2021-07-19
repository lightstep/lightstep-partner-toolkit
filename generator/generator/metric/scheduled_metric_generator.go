package metric

import (
	"log"
	"time"

	"github.com/smithclay/synthetic-load-generator-go/emitter"
	"github.com/smithclay/synthetic-load-generator-go/topology"
)

type ScheduledMetricGenerator struct {
	metrics        []topology.Metric
	service        string
	metricsPerHour int
	seed           int64
	Emitter        emitter.MetricEmitter
	ticker         *time.Ticker
	metricCount    int
	closed         chan struct{}
}

type ScheduledMetricGeneratorOption func(*ScheduledMetricGenerator)

func WithSeed(seed int64) ScheduledMetricGeneratorOption {
	return func(h *ScheduledMetricGenerator) {
		h.seed = seed
	}
}

func WithEmitter(e emitter.MetricEmitter) ScheduledMetricGeneratorOption {
	return func(h *ScheduledMetricGenerator) {
		h.Emitter = e
	}
}

func WithMetricsPerHour(metricsPerHour int) ScheduledMetricGeneratorOption {
	return func(h *ScheduledMetricGenerator) {
		h.metricsPerHour = metricsPerHour
	}
}

func NewScheduledMetricGenerator(metrics []topology.Metric, service string, opts ...ScheduledMetricGeneratorOption) *ScheduledMetricGenerator {
	const (
		defaultMetricsPerHour = 720
		defaultSeed           = 42
	)

	stg := &ScheduledMetricGenerator{
		metrics:        metrics,
		service:        service,
		seed:           defaultSeed,
		metricsPerHour: defaultMetricsPerHour,
		Emitter:        emitter.NewOpenTelemetryStdoutEmitter(),
	}

	for _, opt := range opts {
		opt(stg)
	}

	return stg
}

func (stg *ScheduledMetricGenerator) emitOneMetric() {
	stg.Emitter.EmitMetric(stg.metrics, stg.service)
	stg.metricCount = stg.metricCount + 1
}

func (stg *ScheduledMetricGenerator) Start() {
	log.Printf("Starting metric generation for service %s %d metrics/hr...",
		stg.service, stg.metricsPerHour)
	stg.ticker = time.NewTicker(time.Duration(360000/stg.metricsPerHour) * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-stg.ticker.C:
				stg.emitOneMetric()
			}
		}
	}()
}

func (stg *ScheduledMetricGenerator) Shutdown() {
	log.Printf("\nSummary: Emitted %v metrics for %v", stg.metricCount, stg.service)
	stg.Emitter.Close()
	stg.ticker.Stop()
}
