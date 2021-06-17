package generator

import (
	"fmt"
	"github.com/smithclay/synthetic-load-generator-go/topology"
	"math/rand"
)
import "github.com/smithclay/synthetic-load-generator-go/trace"

type TraceGenerator struct {
	topology *topology.Topology
	trace *trace.Trace
	sequenceNumber int
	random *rand.Rand
	tagNameGenerator topology.Generator
}

func NewTraceGenerator(t *topology.Topology, seed int64) *TraceGenerator {
	r := rand.New(rand.NewSource(seed))
	r.Seed(seed)

	tg := &TraceGenerator{
		topology: t,
		trace: trace.NewTrace(),
		random: r,
	}
	return tg
}

func (g *TraceGenerator) Generate(rootServiceName string, rootRouteName string, startTimeMicros int64) *trace.Trace {
	rootService := g.topology.GetServiceTier(rootServiceName)
	rootSpan := g.createSpanForServiceRouteCall(rootService, rootRouteName, startTimeMicros)
	g.trace.RootSpan = rootSpan
	g.trace.AddRefs()
	return g.trace
}

func Max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

func (g *TraceGenerator) createSpanForServiceRouteCall(serviceTier *topology.ServiceTier, routeName string, startTimeMicros int64) *trace.Span {
	serviceTier.Random = g.random

	instanceName := serviceTier.GetRandomInstance()
	route := serviceTier.GetRoute(routeName)

	service := trace.Service{
		ServiceName:  serviceTier.ServiceName,
		InstanceName: instanceName,
	}

	span := trace.NewSpan(service, route.Route, startTimeMicros)
	span.AddTagString("load_generator.seq_num", fmt.Sprintf("%v", g.sequenceNumber))

	tagSet := serviceTier.GetTagSet(routeName)
	for _, ts := range tagSet {
		for k, v := range ts.Tags {
			span.AddTagString(k, fmt.Sprintf("%v", v))
		}
		for _, tg := range ts.TagGenerators {
			tg.Random = g.random
			for k, v := range tg.GenerateTags() {
				span.AddTagString(k, fmt.Sprintf("%v", v))
			}
		}
	}
	maxEndTime := startTimeMicros
	for s, r := range route.DownstreamCalls {
		childStartTimeMicros := startTimeMicros + (g.random.Int63n(route.MaxLatencyMillis * 1000000))
		childSvc := g.topology.GetServiceTier(s)
		childSpan := g.createSpanForServiceRouteCall(childSvc, r, childStartTimeMicros)
		ref := trace.Reference{
			FromSpanId: span.ID,
			ToSpanId:   childSpan.ID,
			RefType:    trace.CHILD_OF,
		}
		childSpan.AddRef(ref)
		maxEndTime = Max(maxEndTime, childStartTimeMicros)
	}
	ownDuration := g.random.Int63n(route.MaxLatencyMillis * 1000000)
	span.EndTimeMicros = maxEndTime + ownDuration
	g.trace.AddSpan(span)
	g.sequenceNumber = g.sequenceNumber + 1
	return span
}
