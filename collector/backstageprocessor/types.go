package backstageprocessor

type BackstageEntityResponse []struct {
	Metadata struct {
		Namespace   string                 `json:"namespace"`
		Annotations map[string]interface{} `json:"annotations"`
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Tags        []string               `json:"tags"`
		UID         string                 `json:"uid"`
		Etag        string                 `json:"etag"`
		Generation  int                    `json:"generation"`
	} `json:"metadata"`
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Spec       struct {
		Type      string `json:"type"`
		Lifecycle string `json:"lifecycle"`
		Owner     string `json:"owner"`
	} `json:"spec"`
	Relations []struct {
		Target struct {
			Kind      string `json:"kind"`
			Namespace string `json:"namespace"`
			Name      string `json:"name"`
		} `json:"target"`
		Type string `json:"type"`
	} `json:"relations"`
}
