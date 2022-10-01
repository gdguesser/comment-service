package main

import (
	"context"
	"log"
	"os"

	"github.com/gdguesser/comment-service/internal/comment"
	"github.com/gdguesser/comment-service/internal/db"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"google.golang.org/grpc/credentials"

	transportHttp "github.com/gdguesser/comment-service/internal/transport/http"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

func initTracer() func(context.Context) error {

	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if len(insecure) > 0 {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		log.Fatal(err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Print("Could not set resources: ", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}

// Run - is going to be responsible for the instantiation and startup of our go application
func Run() error {
	db, err := db.NewDatabase()
	if err != nil {
		log.Println("failed to connect to the database")
		return err
	}
	if err := db.MigrateDB(); err != nil {
		log.Println("failed to migrate database")
		return err
	}
	cmtService := comment.NewService(db)

	httpHandler := transportHttp.NewHandler(cmtService)
	log.Printf("starting up our application on %v", httpHandler.Server.Addr)
	if err := httpHandler.Serve(); err != nil {
		return err
	}

	return nil
}

func main() {
	cleanup := initTracer()
	defer cleanup(context.Background())

	if err := Run(); err != nil {
		log.Printf("Run failed: %s\n", err)
	}
}
