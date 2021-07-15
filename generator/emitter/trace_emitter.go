package emitter

import "github.com/smithclay/synthetic-load-generator-go/trace"

type TraceEmitter interface {
	EmitTrace(t *trace.Trace)
	Close()
}