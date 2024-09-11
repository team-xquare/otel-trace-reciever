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
		traces := telemetry.ConvertResourceSpansToTraces(resourceSpans)
		if traces != nil {
			err := s.traceService.ProcessTrace(ctx, traces)
			if err != nil {
				fmt.Sprintf("Failed to process trace: %v", err)
				return nil, status.Errorf(codes.Internal, "Failed to process trace: %v", err)
			}
		}
	}

	return &collectorpb.ExportTraceServiceResponse{}, nil
}
