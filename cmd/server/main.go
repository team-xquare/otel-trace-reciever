package main

import (
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
	cfg, err := config.Load()
	repo, err := repository.NewMongoRepository(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	traceService := service.NewTraceService(repo)
	traceServer := server.NewTraceServer(traceService)

	grpcServer := grpc.NewServer()
	collectorpb.RegisterTraceServiceServer(grpcServer, traceServer)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", "0.0.0.0:4317")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("gRPC server listening on 0.0.0.0:4317")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
