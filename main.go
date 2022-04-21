package main

import (
	"auth_api/config"
	"auth_api/server"
	"auth_api/telemetry"
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {
	log.Printf("msg=\"setting up application...\", version=\"%s\", app=\"auth_api\", level=\"info\"", config.Version)
	s := server.New("8840")

	tp, err := telemetry.NewTracerProvider("http://192.168.1.50:14268/api/traces")

	if err != nil {
		log.Printf("msg=\"failed to create tracer provider\", err=\"%s\", app=\"auth_api\", level=\"error\"", err)
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("msg=\"failed to shutdown tracer\", err=\"%s\", app=\"auth_api\", level=\"error\"", err)
		}
	}()

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	log.Printf("msg=\"application started\", version=\"%s\", app=\"auth_api\", level=\"info\"", config.Version)
	s.Run()
}
