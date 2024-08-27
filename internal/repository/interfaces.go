package repository

import (
	"context"
	"otel-trace-reciever/internal/models"
)

type Repository interface {
	SaveTrace(ctx context.Context, trace *models.Trace) error
}
