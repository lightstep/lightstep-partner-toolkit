package emitter

import "github.com/smithclay/synthetic-load-generator-go/trace"

type TraceEmitter interface {
	Emit(t *trace.Trace)
	Close()
}