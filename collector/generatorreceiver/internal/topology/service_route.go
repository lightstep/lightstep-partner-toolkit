package topology

type ServiceRoute struct {
	Route string `json:"route" yaml:"route"`
	DownstreamCalls map[string]string `json:"downstreamCalls,omitempty" yaml:"downstreamCalls,omitempty"`
	MaxLatencyMillis int64 `json:"maxLatencyMillis" yaml:"maxLatencyMillis"`
	TagSets []TagSet `json:"tagSets" yaml:"tagSets"`
	ResourceAttributeSets []ResourceAttributeSet `json:"resourceAttrSets" yaml:"resourceAttrSets"`
	FlagSet string `json:"flag_set" yaml:"flag_set"`
	FlagUnset string `json:"flag_unset" yaml:"flag_unset"`
}
