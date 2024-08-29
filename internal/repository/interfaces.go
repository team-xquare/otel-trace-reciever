package repository

import (
	"context"
	"otel-trace-reciever/internal/models"
)

type Repository interface {
	SaveTraces(ctx context.Context, traces []*models.Trace) error
}
