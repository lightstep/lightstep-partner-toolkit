package emitter

import "github.com/smithclay/synthetic-load-generator-go/trace"

type MetricEmitter interface {
	EmitMetric(t *trace.Trace)
	Close()
}