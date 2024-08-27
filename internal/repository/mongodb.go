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

func (r *MongoRepository) SaveTrace(ctx context.Context, trace *models.Trace) error {
	_, err := r.traceCollection.InsertOne(ctx, trace)
	return err
}
