package streamreceiver

import "time"

type LightstepStreamResponse struct {
	Data LightstepStreamSlice `json:"data"`
}

type LightstepStreamSlice struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		PointsCount int `json:"points-count"`
		TimeWindows []struct {
			OldestTime   time.Time `json:"oldest-time"`
			YoungestTime time.Time `json:"youngest-time"`
		} `json:"time-windows"`
		TimeRange struct {
			OldestTime   time.Time `json:"oldest-time"`
			YoungestTime time.Time `json:"youngest-time"`
		} `json:"time-range"`
		ResolutionMs int `json:"resolution-ms"`
		Exemplars    []struct {
			DurationMicros int     `json:"duration_micros"`
			OldestMicros   int64   `json:"oldest_micros"`
			YoungestMicros int64   `json:"youngest_micros"`
			SpanGUID       string  `json:"span_guid"`
			HasError       bool    `json:"has_error"`
			Percentile     float64 `json:"percentile"`
			SpanName       string  `json:"span_name"`
			TraceGUID      string  `json:"trace_guid"`
			TraceHandle    string  `json:"trace_handle"`
		} `json:"exemplars"`
	} `json:"attributes"`
}

type LightstepTraceResponse struct {
	Data []LightstepTrace `json:"data""`
}

type LightstepTrace struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		StartTimeMicros int64 `json:"start-time-micros"`
		EndTimeMicros   int64 `json:"end-time-micros"`
		Spans           []struct {
			SpanName        string                 `json:"span-name"`
			SpanID          string                 `json:"span-id"`
			IsError         bool                   `json:"is-error"`
			StartTimeMicros int64                  `json:"start-time-micros"`
			EndTimeMicros   int64                  `json:"end-time-micros"`
			TraceID         string                 `json:"trace-id"`
			ReporterID      string                 `json:"reporter-id"`
			Tags            map[string]interface{} `json:"tags,omitempty"`
		} `json:"spans"`
	} `json:"attributes"`
	Relationships struct {
		Reporters []struct {
			ReporterID string                 `json:"reporter-id"`
			Attributes map[string]interface{} `json:"attributes,omitempty"`
		} `json:"reporters"`
	} `json:"relationships"`
}
