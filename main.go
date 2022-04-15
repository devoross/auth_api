package main

import (
	"auth_api/server"
	"auth_api/telemetry"
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Test struct {
	Test string `json:"test"`
}

func main() {
	log.Println("setting up application...")
	s := server.New("8840")

	tp, err := telemetry.NewTracerProvider("http://192.168.1.50:14268/api/traces")

	if err != nil {
		log.Printf("msg=\"failed to create tracer provider\", err=\"%s\"", err)
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("msg=\"failed to shutdown tracer\", err=\"%s\"", err)
		}
	}()

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	log.Println("application started...")
	s.Run()
}
