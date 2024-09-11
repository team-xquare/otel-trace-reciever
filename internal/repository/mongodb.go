package repository

import (
	"context"
	"fmt"
	"log"
	"otel-trace-reciever/internal/models"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
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

func sanitizeFieldName(name string) string {
	return strings.ReplaceAll(name, ".", "_")
}

func sanitizeDocument(doc bson.M) bson.M {
	result := make(bson.M)
	for k, v := range doc {
		newKey := sanitizeFieldName(k)
		switch val := v.(type) {
		case bson.M:
			result[newKey] = sanitizeDocument(val)
		case []interface{}:
			result[newKey] = sanitizeArray(val)
		default:
			result[newKey] = val
		}
	}
	return result
}

func sanitizeArray(arr []interface{}) []interface{} {
	result := make([]interface{}, len(arr))
	for i, v := range arr {
		if doc, ok := v.(bson.M); ok {
			result[i] = sanitizeDocument(doc)
		} else {
			result[i] = v
		}
	}
	return result
}

func (r *MongoRepository) SaveTraces(ctx context.Context, traces []*models.Trace) error {
	if len(traces) == 0 {
		return nil
	}

	documents := make([]interface{}, len(traces))
	for i, trace := range traces {
		// Convert trace to BSON and sanitize field names
		bsonDoc, err := bson.Marshal(trace)
		if err != nil {
			return fmt.Errorf("failed to marshal trace: %v", err)
		}
		var doc bson.M
		err = bson.Unmarshal(bsonDoc, &doc)
		if err != nil {
			return fmt.Errorf("failed to unmarshal trace: %v", err)
		}
		sanitizedDoc := sanitizeDocument(doc)
		documents[i] = sanitizedDoc
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
