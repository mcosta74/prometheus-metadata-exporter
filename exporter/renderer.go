package exporter

import (
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"sort"
	txttemplate "text/template"
)

type renderer interface {
	Render(metaInfo data) error
}

type consoleRenderer struct{}

func (r *consoleRenderer) Render(metaInfo data) error {
	keys := make([]string, 0, len(metaInfo))
	for k, _ := range metaInfo {
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
	keys := make([]string, 0, len(metaInfo))
	for k, _ := range metaInfo {
		keys = append(keys, k)
	}
	sort.Strings(keys)

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
	keys := make([]string, 0, len(metaInfo))
	for k, _ := range metaInfo {
		keys = append(keys, k)
	}
	sort.Strings(keys)

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

type htmlRenderer struct {
	fs fs.FS
}

func (r *htmlRenderer) Render(metaInfo data) error {
	keys := make([]string, 0, len(metaInfo))
	for k, _ := range metaInfo {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	t, err := template.ParseFS(r.fs, "templates/output.html")
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

type mdRenderer struct {
	fs fs.FS
}

func (r *mdRenderer) Render(metaInfo data) error {
	keys := make([]string, 0, len(metaInfo))
	for k, _ := range metaInfo {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	t, err := txttemplate.ParseFS(r.fs, "templates/output.md")
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
