module github.com/kit494way/k8s-hpa-opentelemetry-cloud-monitoring/app

go 1.15

require (
	cloud.google.com/go v0.76.0
	go.opentelemetry.io/contrib/detectors/gcp v0.16.0
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/exporters/otlp v0.16.0 // indirect
	go.opentelemetry.io/otel/sdk v0.16.0
)
