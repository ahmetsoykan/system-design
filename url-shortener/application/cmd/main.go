package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"url-shortener/cmd/handlers"

	"github.com/kelseyhightower/envconfig"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {

	var s handlers.Config
	err := envconfig.Process("app", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	// initialize the tracer provider and exporter
	err = initTracer()
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}

	log.Printf(fmt.Sprintf("main : app started on localhost:%s", s.Port))

	api := http.Server{
		Addr:    fmt.Sprintf("localhost:%s", s.Port),
		Handler: otelhttp.NewHandler(handlers.NewServer(s).Router, "http-server"),
	}

	if err := api.ListenAndServe(); err != nil {
		log.Fatal("error: %s", err)
	}
}

// AWS X-Ray Distributed Tracing
// docker run --rm -d -p 4317:4317 -p 55680:55680 -p 8889:8888 -e "AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}" -e "AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}" -e AWS_REGION=eu-west-1 --name awscollector public.ecr.aws/aws-observability/aws-otel-collector:latest
func initTracer() error {

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint("localhost:4317"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "url-shortener"),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Printf("could not set resources: ", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return nil
}
