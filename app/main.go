package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	// "google.golang.org/grpc"
	gcpdetector "go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/metric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
)

func initProvider() func() {
	ctx := context.Background()

	collectorEndpoint := os.Getenv("OTEL_COLLECTOR_ENDPOINT")
	if collectorEndpoint == "" {
		log.Fatal("Failed to get environment variable OTEL_COLLECTOR_ENDPOINT.")
	}

	driver := otlpgrpc.NewDriver(
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint(collectorEndpoint),
		// otlpgrpc.WithDialOption(grpc.WithBlock()), // useful for testing
	)
	exporter, err := otlp.NewExporter(ctx, driver)
	if err != nil {
		log.Fatalln("Failed to create exporter:", err)
	}

	res, err := resource.New(
		ctx,
		resource.WithDetectors(&gcpdetector.GKE{}),
	)
	if err != nil {
		log.Fatalln("could not initialize resource:", err)
	}

	cont := controller.New(
		processor.New(
			// simple.NewWithInexpensiveDistribution(),
			simple.NewWithExactDistribution(),
			exporter,
		),
		controller.WithPusher(exporter),
		controller.WithCollectPeriod(10*time.Second),
		controller.WithResource(res),
	)

	otel.SetMeterProvider(cont.MeterProvider())
	if err := cont.Start(ctx); err != nil {
		log.Fatalln("Failed to start controller:", err)
	}

	return func() {
		// Push any last metric events to the exporter.
		if err := cont.Stop(ctx); err != nil {
			log.Fatalln("Failed to stop controller:", err)
		}
	}
}

func main() {
	var metricValue float64
	flag.Float64Var(&metricValue, "metric-value", 0, "metric value")
	flag.Parse()

	shutdown := initProvider()
	defer shutdown()

	meter := otel.Meter("test-meter")

	metric.Must(meter).NewFloat64ValueObserver("metric-a", func(_ context.Context, result metric.Float64ObserverResult) {
		result.Observe(metricValue)
	})

	for {
		// do something

		<-time.After(time.Second)
	}
}
