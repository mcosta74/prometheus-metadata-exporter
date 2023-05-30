# prometheus-metadata-exporter

![build](https://github.com/mcosta74/prometheus-metadata-exporter/actions/workflows/build.yml/badge.svg)
![GitHub](https://img.shields.io/github/license/mcosta74/prometheus-metadata-exporter)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/mcosta74/prometheus-metadata-exporter)

Utility to get information about all metrics metadata available on a Prometheus instance

# Usage

Run the tool specifying the Prometheus instance URL (default "http://localhost:9090") and the output format (default "text")

```sh
./prometheus-metadata-exporter -prometheus.url http://<some-server>:<some-port> -format json
```

The tool will print the results on the standard output; if you want redirect to a file you can use shell redirection

```sh
./prometheus-metadata-exporter -prometheus.url http://<some-server>:<some-port> -format json > output.json
```

## Configuration parameter

| Flag Name         | Description                          | Default Value           |
| ----------------- | ------------------------------------ | ----------------------- |
| `-prometheus.url` | Base URL for the Prometheus instance | `http://localhost:9090` |
| `-format`         | Output format                        | `text`                  |

## Output formats

| Option | Description            |
| ------ | ---------------------- |
| `text` | Text Format            |
| `csv`  | CSV Format             |
| `html` | HTML Table Format      |
| `md`   | Markdown Table Format  |
| `json` | JSON (indented) Format |