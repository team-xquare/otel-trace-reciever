package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"otel-trace-reciever/internal/config"
	"otel-trace-reciever/internal/repository"
	"otel-trace-reciever/internal/server"
	"google.golang.org/grpc/reflection"
	"otel-trace-reciever/internal/service"

	collectorpb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
)

func main() {
	grpcHost := "0.0.0.0"
	grpcPort := 4317
	grpcAddress := fmt.Sprintf("%s:%d", grpcHost, grpcPort)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	repo, err := repository.NewMongoRepository(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	traceService := service.NewTraceService(repo)
	traceServer := server.NewTraceServer(traceService)

	grpcServer := grpc.NewServer()
	collectorpb.RegisterTraceServiceServer(grpcServer, traceServer)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on %s", grpcAddress)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}