package server

import (
	"context"
	"fmt"
	"otel-trace-reciever/internal/service"
	"otel-trace-reciever/pkg/telemetry"

	collectorpb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TraceServer struct {
	collectorpb.UnimplementedTraceServiceServer
	traceService *service.TraceService
}

func NewTraceServer(traceService *service.TraceService) *TraceServer {
	return &TraceServer{
		traceService: traceService,
	}
}

func (s *TraceServer) Export(ctx context.Context, req *collectorpb.ExportTraceServiceRequest) (*collectorpb.ExportTraceServiceResponse, error) {
	for _, resourceSpans := range req.ResourceSpans {
		spans := telemetry.ConvertResourceSpansToSpans(resourceSpans)

		err := s.traceService.ProcessSpan(ctx, &spans)
		if err != nil {
			fmt.Printf("Failed to process span: %v\n", err)
			return nil, status.Errorf(codes.Internal, "Failed to process span: %v", err)
		}
	}

	return &collectorpb.ExportTraceServiceResponse{}, nil
}
