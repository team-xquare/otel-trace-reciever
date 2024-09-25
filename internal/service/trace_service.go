package service

import (
	"context"
	"otel-trace-reciever/internal/models"
	"otel-trace-reciever/internal/repository"
)

type TraceService struct {
	repo repository.Repository
}

func NewTraceService(repo repository.Repository) *TraceService {
	return &TraceService{repo: repo}
}

func (s *TraceService) ProcessTrace(ctx context.Context, traces []*models.Trace) error {
	return s.repo.SaveTraces(ctx, traces)
}

func (s *TraceService) ProcessSpan(ctx context.Context, spans *[]models.Span) error {
	return s.repo.SaveSpans(ctx, spans)
}
