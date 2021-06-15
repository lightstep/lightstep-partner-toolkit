package trace

type RefType int

const (
	CHILD_OF RefType = 0
)

type KeyValue struct {
	Key string
	Value string
}

type Reference struct {
	FromSpanId [8]byte
	ToSpanId [8]byte
	RefType RefType
}

type Service struct {
	ServiceName string
	InstanceName string
}

type Trace struct  {
	RootSpan *Span
	Spans []*Span
	SpanIdToSpan map[[8]byte]*Span
	SpanIdToOutgoingRefs map[[8]byte][]Reference
}

func NewTrace() *Trace {
	return &Trace{
		Spans: make([]*Span, 0),
		SpanIdToSpan: make(map[[8]byte]*Span),
		SpanIdToOutgoingRefs: make(map[[8]byte][]Reference),
	}
}

func (t *Trace) AddSpan(span *Span) *[]*Span {
	t.Spans = append(t.Spans, span)
	t.SpanIdToSpan[span.ID] = span
	return &t.Spans
}

func (t *Trace) AddRefs() {
	for _, span := range t.Spans {
		for _, ref := range span.Refs {
			if _, ok := t.SpanIdToOutgoingRefs[ref.FromSpanId]; !ok {
				t.SpanIdToOutgoingRefs[ref.FromSpanId] = make([]Reference, 0)
			}
			t.SpanIdToOutgoingRefs[ref.FromSpanId] = append(t.SpanIdToOutgoingRefs[ref.FromSpanId], ref)
		}
	}
}