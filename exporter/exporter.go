package exporter

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
)

type Exporter interface {
	Export(baseUrl string) error
}

func NewExporter(format string, templatesFS fs.FS) Exporter {
	return &exporter{
		r: newRenderer(format, templatesFS),
	}
}

type exporter struct {
	r renderer
}

func (exp *exporter) Export(baseUrl string) error {
	data, err := getMetadata(baseUrl)
	if err != nil {
		return err
	}

	return exp.r.Render(data)
}

func getMetadata(baseUrl string) (data, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/metadata", baseUrl))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status: %s", resp.Status)
	}

	var rData metadataResponse
	if err := json.NewDecoder(resp.Body).Decode(&rData); err != nil {
		return nil, err
	}
	return rData.Data, nil
}
