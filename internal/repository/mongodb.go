package repository

import (
	"context"
	"otel-trace-reciever/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	client          *mongo.Client
	traceCollection *mongo.Collection
}

func NewMongoRepository(uri string) (*MongoRepository, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &MongoRepository{
		client:          client,
		traceCollection: client.Database("tracing").Collection("traces"),
	}, nil
}

func (r *MongoRepository) SaveTraces(ctx context.Context, traces []*models.Trace) error {
	if len(traces) == 0 {
		return nil
	}

	documents := make([]interface{}, len(traces))
	for i, trace := range traces {
		documents[i] = trace
	}

	_, err := r.traceCollection.InsertMany(ctx, documents)
	return err
}
