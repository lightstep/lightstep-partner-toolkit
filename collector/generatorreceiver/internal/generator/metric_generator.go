package generator

import (
	"go.opentelemetry.io/collector/model/pdata"
	"math"
	"math/rand"
	"time"
)

type MetricGenerator struct {
	metricCount int
	random *rand.Rand
}

func NewMetricGenerator(seed int64) *MetricGenerator {
	r := rand.New(rand.NewSource(seed))
	r.Seed(seed)
	return &MetricGenerator{
		metricCount: 0,
		random: r,
	}
}

func (g *MetricGenerator) Generate(metricName string, metricType string, serviceName string) pdata.Metrics {
	metrics := pdata.NewMetrics()

	rms := metrics.ResourceMetrics().AppendEmpty()
	rms.Resource().Attributes().InsertString("service.name", serviceName)

	m := rms.InstrumentationLibraryMetrics().AppendEmpty().Metrics().AppendEmpty()
	m.SetName(metricName)
	if metricType == "Gauge" {
		m.SetDataType(pdata.MetricDataTypeGauge)
		dp := m.Gauge().DataPoints().AppendEmpty()
		dp.SetTimestamp(pdata.NewTimestampFromTime(time.Now()))
		dp.SetDoubleVal((math.Sin(.01 * float64(g.metricCount)) + 1) * 100)
	} else if metricType == "Sum" {
		m.SetDataType(pdata.MetricDataTypeSum)
		dp := m.Sum().DataPoints().AppendEmpty()
		dp.SetTimestamp(pdata.NewTimestampFromTime(time.Now()))
		dp.SetDoubleVal((math.Sin(.01 * float64(g.metricCount)) + 1) * 100)
	}

	g.metricCount = g.metricCount + 1
	return metrics
}