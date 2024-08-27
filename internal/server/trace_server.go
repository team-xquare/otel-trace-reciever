package server

import (
	"context"
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
		trace := telemetry.ConvertResourceSpansToTrace(resourceSpans)
		if trace != nil {
			err := s.traceService.ProcessTrace(ctx, trace)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "Failed to process trace: %v", err)
			}
		}
	}

	return &collectorpb.ExportTraceServiceResponse{}, nil
}
