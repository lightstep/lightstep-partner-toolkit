package topology

import (
	"github.com/lightstep/lightstep-partner-sdk/collector/generatorreceiver/internal/flags"
	"github.com/lightstep/lightstep-partner-sdk/collector/generatorreceiver/internal/incidents"
)

type File struct {
	Topology   *Topology    `json:"topology" yaml:"topology"`
	Flags      []flags.Flag `json:"flags" yaml:"flags"`
	RootRoutes []struct {
		Service       string `json:"service" yaml:"service"`
		Route         string `json:"route" yaml:"route"`
		TracesPerHour int    `json:"tracesPerHour" yaml:"tracesPerHour"`
	} `json:"rootRoutes" yaml:"rootRoutes"`
	Incidents []incidents.Incident `json:"incidents" yaml:"incidents"`
}
