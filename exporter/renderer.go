package exporter

import (
	"encoding/json"
	html "html/template"
	"io"
	"io/fs"
	"os"
	"sort"
	txt "text/template"
)

type renderer interface {
	Render(metaInfo data) error
}

func newRenderer(format string, templatesFS fs.FS) renderer {
	var r renderer
	switch format {
	case "csv":
		r = &templateRenderer[*txt.Template]{
			fs:           templatesFS,
			templateName: "templates/csv.tmpl",
			parseFunc:    txt.ParseFS,
		}

	case "html":
		r = &templateRenderer[*html.Template]{
			fs:           templatesFS,
			templateName: "templates/html.tmpl",
			parseFunc:    html.ParseFS,
		}

	case "md":
		r = &templateRenderer[*txt.Template]{
			fs:           templatesFS,
			templateName: "templates/md.tmpl",
			parseFunc:    txt.ParseFS,
		}

	case "json":
		r = &jsonRenderer{}

	default:
		r = &templateRenderer[*txt.Template]{
			fs:           templatesFS,
			templateName: "templates/text.tmpl",
			parseFunc:    txt.ParseFS,
		}
	}
	return r
}

type jsonRenderer struct{}

type record struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Help string `json:"help,omitempty"`
	Unit string `json:"unit,omitempty"`
}

func (r *jsonRenderer) Render(metaInfo data) error {
	keys := sortedKeys(metaInfo)

	records := make([]record, 0)
	for _, k := range keys {
		for _, item := range metaInfo[k] {
			records = append(records, record{k, item.Type, item.Help, item.Unit})
		}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

type Executor interface {
	Execute(io.Writer, any) error
}

type templateRenderer[T Executor] struct {
	fs           fs.FS
	templateName string
	parseFunc    func(fs.FS, ...string) (T, error)
}

func (r *templateRenderer[T]) Render(metaInfo data) error {
	keys := sortedKeys(metaInfo)

	t, err := r.parseFunc(r.fs, r.templateName)
	if err != nil {
		return err
	}

	records := make([]record, 0)
	for _, k := range keys {
		for _, item := range metaInfo[k] {
			records = append(records, record{k, item.Type, item.Help, item.Unit})
		}
	}
	return t.Execute(os.Stdout, records)
}

func sortedKeys(metaInfo data) []string {
	keys := make([]string, 0, len(metaInfo))
	for k, _ := range metaInfo {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
