package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"github.com/mcosta74/prometheus-metadata-exporter/exporter"
)

//go:embed templates/*
var templatesDir embed.FS

var (
	prometheusUrl *string
	format        *string
)

func main() {
	fs := flag.NewFlagSet("prometheus-metadata-exporter", flag.ExitOnError)
	initFlags(fs)
	fs.Parse(os.Args[1:])

	exp := exporter.NewExporter(*format, templatesDir)
	err := exp.Export(*prometheusUrl)
	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}
}

func initFlags(fs *flag.FlagSet) {
	prometheusUrl = fs.String("prometheus.url", "http://localhost:9090", "URL of the Prometheus server to get information from.")
	format = fs.String("format", "text", "Output format (text, csv, html, md, json)")
}
