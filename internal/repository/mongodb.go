package repository

import (
	"context"
	"fmt"
	"log"
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
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB")

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

	result, err := r.traceCollection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to insert traces: %v", err)
	}

	log.Printf("Successfully inserted %d traces", len(result.InsertedIDs))
	return nil
}

func (r *MongoRepository) Close(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}
