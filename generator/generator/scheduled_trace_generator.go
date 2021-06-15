package generator

import (
	"github.com/smithclay/synthetic-load-generator-go/emitter"
	"github.com/smithclay/synthetic-load-generator-go/topology"
	"log"
	"time"
)

type ScheduledTraceGenerator struct {
	topology *topology.Topology
	route string
	service string
	tracesPerHour int
	Emitter emitter.TraceEmitter
	ticker *time.Ticker
	traceGen *TraceGenerator
	traceCount int
	closed chan struct{}
}

type ScheduledTraceGeneratorOption func(*ScheduledTraceGenerator)

func WithSeed(seed int64) ScheduledTraceGeneratorOption {
	return func(h *ScheduledTraceGenerator) {
		h.traceGen = NewTraceGenerator(h.topology, seed)
	}
}

func WithGrpc(url string) ScheduledTraceGeneratorOption {
	return func(h *ScheduledTraceGenerator) {
		h.Emitter = emitter.NewOpenTelemetryGrpcEmitter(url)
	}
}

func WithTracesPerHour(tracesPerHour int) ScheduledTraceGeneratorOption {
	return func(h *ScheduledTraceGenerator) {
		h.tracesPerHour = tracesPerHour
	}
}

func NewScheduledTraceGenerator(topo *topology.Topology, route string, service string, opts ...ScheduledTraceGeneratorOption) *ScheduledTraceGenerator {
	const (
		defaultTracesPerHour = 720
		defaultSeed  = 42
	)

	stg := &ScheduledTraceGenerator{
		topology: topo,
		route: route,
		service: service,
		tracesPerHour: defaultTracesPerHour,
		traceGen: NewTraceGenerator(topo, defaultSeed),
		Emitter: emitter.NewOpenTelemetryStdoutEmitter(),
	}

	for _, opt := range opts {
		opt(stg)
	}

	return stg
}

func (stg *ScheduledTraceGenerator) emitOneTrace() {
	t := stg.traceGen.Generate(stg.service, stg.route, time.Now().UnixNano())
	stg.Emitter.Emit(t)
	stg.traceCount = stg.traceCount + 1
}

func (stg *ScheduledTraceGenerator) Start() {
	log.Printf("Starting trace generation for service %s, route %s, %d traces/hr...",
		stg.service, stg.route, stg.tracesPerHour)
	stg.ticker = time.NewTicker(time.Duration(360000/stg.tracesPerHour) * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-stg.ticker.C:
				stg.emitOneTrace()
			}
		}
	}()
}

func (stg *ScheduledTraceGenerator) Shutdown() {
	log.Printf("\nSummary: Emitted %v traces for %v%v", stg.traceCount, stg.service, stg.route)
	stg.Emitter.Close()
	stg.ticker.Stop()
}