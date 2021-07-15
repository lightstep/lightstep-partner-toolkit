package emitter

import "github.com/smithclay/synthetic-load-generator-go/topology"

type MetricEmitter interface {
	EmitMetric(metrics []topology.Metric, service string)
	Close()
}