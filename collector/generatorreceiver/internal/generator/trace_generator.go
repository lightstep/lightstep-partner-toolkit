package generator

import (
	"fmt"
	"github.com/lightstep/lightstep-partner-sdk/collector/generatorreceiver/internal/topology"
	"go.opentelemetry.io/collector/model/pdata"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"math/rand"
	"sync"
	"time"
)

type TraceGenerator struct {
	topology *topology.Topology
	service string
	route string
	sequenceNumber int
	random *rand.Rand
	sync.Mutex
	tagNameGenerator topology.Generator
}

func NewTraceGenerator(t *topology.Topology, seed int64, service string, route string) *TraceGenerator {
	r := rand.New(rand.NewSource(seed))
	r.Seed(seed)

	tg := &TraceGenerator{
		topology: t,
		random: r,
		service: service,
		route: route,
	}
	return tg
}

func (g *TraceGenerator) genTraceId() pdata.TraceID {
	g.Lock()
	defer g.Unlock()
	traceIdBytes := make([]byte, 16)
	g.random.Read(traceIdBytes)
	var traceId [16]byte
	copy(traceId[:], traceIdBytes)
	return pdata.NewTraceID(traceId)
}

func (g *TraceGenerator) genSpanId() pdata.SpanID {
	g.Lock()
	defer g.Unlock()
	traceIdBytes := make([]byte, 16)
	g.random.Read(traceIdBytes)
	var traceId [8]byte
	copy(traceId[:], traceIdBytes)
	return pdata.NewSpanID(traceId)
}

func (g *TraceGenerator) Generate(startTimeMicros int64) *pdata.Traces {
	rootService := g.topology.GetServiceTier(g.service)
	traces := pdata.NewTraces()
	g.createSpanForServiceRouteCall(&traces, rootService, g.route, startTimeMicros, g.genTraceId(), pdata.NewSpanID([8]byte{0x0}))
	return &traces
}

func (g *TraceGenerator) createSpanForServiceRouteCall(traces *pdata.Traces, serviceTier *topology.ServiceTier, routeName string, startTimeMicros int64, traceId pdata.TraceID, parentSpanId pdata.SpanID) *pdata.Span {
	serviceTier.Random = g.random
	route := serviceTier.GetRoute(routeName)

	rspanSlice := traces.ResourceSpans()
	rspan := rspanSlice.AppendEmpty()

	resource := rspan.Resource()

	resource.Attributes().InsertString(string(semconv.ServiceNameKey), serviceTier.ServiceName)

	resourceAttributeSet := serviceTier.GetResourceAttributeSet(); if resourceAttributeSet != nil {
		for k, v := range resourceAttributeSet.ResourceAttributes {
			resource.Attributes().InsertString(k, fmt.Sprintf("%v", v))
		}
	}


	ils := rspan.InstrumentationLibrarySpans().AppendEmpty()
	spans := ils.Spans()

	span := spans.AppendEmpty()
	newSpanId := g.genSpanId()
	span.SetName(routeName)
	span.SetTraceID(traceId)
	span.SetParentSpanID(parentSpanId)
	span.SetSpanID(newSpanId)
	span.SetKind(pdata.SpanKindServer)
	span.Attributes().InsertString("load_generator.seq_num", fmt.Sprintf("%v", g.sequenceNumber))

	tagSet := serviceTier.GetTagSet(routeName)
	for _, ts := range tagSet {
		for k, v := range ts.Tags {
			span.Attributes().InsertString(k, fmt.Sprintf("%v", v))
		}
		for _, tg := range ts.TagGenerators {
			tg.Random = g.random
			for k, v := range tg.GenerateTags() {
				span.Attributes().InsertString(k, fmt.Sprintf("%v", v))
			}
		}
	}

	maxEndTime := startTimeMicros
	for s, r := range route.DownstreamCalls {
		childStartTimeMicros := startTimeMicros + (g.random.Int63n(route.MaxLatencyMillis * 1000000))
		childSvc := g.topology.GetServiceTier(s)
		g.createSpanForServiceRouteCall(traces, childSvc, r, childStartTimeMicros, traceId, newSpanId)
		maxEndTime = Max(maxEndTime, childStartTimeMicros)
	}
	ownDuration := g.random.Int63n(route.MaxLatencyMillis * 1000000)
	span.SetStartTimestamp(pdata.NewTimestampFromTime(time.Unix(0, startTimeMicros)))
	span.SetEndTimestamp(pdata.NewTimestampFromTime(time.Unix(0, maxEndTime + ownDuration)))
	g.sequenceNumber = g.sequenceNumber + 1
	return &span
}

func Max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}