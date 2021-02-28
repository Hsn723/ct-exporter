# ct-exporter
[![GoDoc](https://godoc.org/github.com/Hsn723/ct-exporter?status.svg)](https://godoc.org/github.com/Hsn723/ct-exporter) [![Go Report Card](https://goreportcard.com/badge/github.com/Hsn723/ct-exporter)](https://goreportcard.com/report/github.com/Hsn723/ct-exporter) ![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/Hsn723/ct-exporter?label=latest%20version) ![Docker Pulls](https://img.shields.io/docker/pulls/hsn723/ct-exporter) ![GitHub](https://img.shields.io/github/license/hsn723/ct-exporter)

Export certificate transparency issuances as Prometheus metrics.

## Sample usage
```yaml
scrape_configs:
  - job_name "ct-exporter"
    metrics_path: /probe
    static_configs:
      - targets:
          - example.com
```
