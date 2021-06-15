package trace

import (
	"math/rand"
)

type Span struct {
	ID [8]byte
	Service Service
	StartTimeMicros int64
	EndTimeMicros int64
	OperationName string
	Refs []Reference
	Tags []KeyValue
}

func NewSpan(service Service, opName string, startTimeMicros int64) *Span {
	var sid [8]byte
	rand.Read(sid[:])
	return &Span{
		ID: sid,
		Service: service,
		Tags: make([]KeyValue, 0),
		Refs: make([]Reference, 0),
		OperationName: opName,
		StartTimeMicros: startTimeMicros,
	}
}

func (s *Span) AddTagString(key string, val string) *[]Reference {
	s.Tags = append(s.Tags, KeyValue{key, val})
	return &s.Refs
}

func (s *Span) AddRef(ref Reference) *[]Reference {
	s.Refs = append(s.Refs, ref)
	return &s.Refs
}