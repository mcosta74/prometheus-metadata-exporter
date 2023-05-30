package exporter

type metadata struct {
	Type string `json:"type,omitempty"`
	Help string `json:"help,omitempty"`
	Unit string `json:"unit,omitempty"`
}

type data map[string][]metadata

type metadataResponse struct {
	Status string `json:"status,omitempty"`
	Data   data   `json:"data,omitempty"`
}
