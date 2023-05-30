package exporter

import (
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"sort"
	txttemplate "text/template"
)

type renderer interface {
	Render(metaInfo data) error
}

func newRenderer(format string, templatesFS fs.FS) renderer {
	var r renderer
	switch format {
	case "text":
		r = &consoleRenderer{}

	case "csv":
		r = &csvRenderer{}

	case "html":
		r = &templateRenderer[*template.Template]{
			fs:           templatesFS,
			templateName: "templates/output.html",
			parseFunc:    template.ParseFS,
		}

	case "md":
		r = &templateRenderer[*txttemplate.Template]{
			fs:           templatesFS,
			templateName: "templates/output.md",
			parseFunc:    txttemplate.ParseFS,
		}

	case "json":
		r = &jsonRenderer{}

	default:
		r = &consoleRenderer{}
	}
	return r
}

type consoleRenderer struct{}

func (r *consoleRenderer) Render(metaInfo data) error {
	keys := make([]string, 0, len(metaInfo))
	for k := range metaInfo {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%q:\n", k)
		for _, item := range metaInfo[k] {
			fmt.Printf("\tType: %q\n", item.Type)
			fmt.Printf("\tHelp: %q\n", item.Help)
			if item.Unit != "" {
				fmt.Printf("\tUnit: %q\n", item.Unit)
			}
		}
	}
	return nil
}

type jsonRenderer struct{}

func (r *jsonRenderer) Render(metaInfo data) error {
	keys := sortedKeys(metaInfo)

	type record struct {
		Name string `json:"name,omitempty"`
		Type string `json:"type,omitempty"`
		Help string `json:"help,omitempty"`
		Unit string `json:"unit,omitempty"`
	}

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

type csvRenderer struct{}

func (r *csvRenderer) Render(metaInfo data) error {
	keys := sortedKeys(metaInfo)

	w := csv.NewWriter(os.Stdout)
	w.Write([]string{"Name, Type, Help, Unit"})

	for _, k := range keys {
		for _, item := range metaInfo[k] {
			w.Write([]string{k, item.Type, item.Help, item.Unit})
		}
	}
	w.Flush()
	return nil
}

type Executer interface {
	Execute(io.Writer, any) error
}

type templateRenderer[T Executer] struct {
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

	type record struct {
		Name string
		Type string
		Help string
		Unit string
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
