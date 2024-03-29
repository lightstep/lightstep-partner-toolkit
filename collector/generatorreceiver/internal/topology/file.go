package topology

type File struct {
	Topology *Topology `json:"topology" yaml:"topology"`
	RootRoutes []struct {
		Service       string `json:"service" yaml:"service"`
		Route         string `json:"route" yaml:"route"`
		TracesPerHour int    `json:"tracesPerHour" yaml:"tracesPerHour"`
	} `json:"rootRoutes" yaml:"rootRoutes"`
}