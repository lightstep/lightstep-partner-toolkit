package emitter

import "github.com/smithclay/synthetic-load-generator-go/trace"

func prePostOrder(trace *trace.Trace, preVisitHandler func(s *trace.Span), postVisitHandler func(s *trace.Span)) {
	prePostOrderWithSpan(trace, trace.RootSpan, preVisitHandler, postVisitHandler)
}

func prePostOrderWithSpan(trace *trace.Trace, span *trace.Span, preVisitHandler func(s *trace.Span), postVisitHandler func(s *trace.Span)) {
	preVisitHandler(span)
	if outgoing, ok := trace.SpanIdToOutgoingRefs[span.ID]; ok {
		for _, ref := range outgoing {
			descendant := trace.SpanIdToSpan[ref.ToSpanId]
			prePostOrderWithSpan(trace, descendant, preVisitHandler, postVisitHandler)
		}
	}
	postVisitHandler(span)
}