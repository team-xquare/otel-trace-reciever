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
		traceCollection: client.Database("tracing").Collection("spans"),
	}, nil
}

func (r *MongoRepository) SaveTraces(ctx context.Context, traces []*models.Trace) error {
	if len(traces) == 0 {
		return nil
	}

	documents := make([]interface{}, len(traces))
	for i, trace := range traces {
		doc := bson.M{
			"traceId":      trace.TraceID,
			"serviceName":  trace.ServiceName,
			"dateNano":     trace.DateNano,
			"durationNano": trace.DurationNano,
			"spans":        convertSpans(trace.Spans),
		}
		documents[i] = doc
	}

	result, err := r.traceCollection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to insert traces: %v", err)
	}

	log.Printf("Successfully inserted %d traces", len(result.InsertedIDs))
	return nil
}

func (r *MongoRepository) SaveSpans(ctx context.Context, spans []models.Span) error {
	if len(spans) == 0 {
		return nil
	}

	documents := make([]interface{}, len(spans))
	for i, span := range spans {
		doc := bson.M{
			"id":                span.ID,
			"traceId":           span.TraceID,
			"spanId":            span.SpanID,
			"parentSpanId":      span.ParentSpanID,
			"name":              span.Name,
			"kind":              span.Kind,
			"startTimeUnixNano": span.StartTimeUnixNano,
			"endTimeUnixNano":   span.EndTimeUnixNano,
			"attributes":        convertAttributes(span.Attributes),
			"events":            convertEvents(span.Events),
			"links":             convertLinks(span.Links),
			"status":            convertStatus(span.Status),
		}
		documents[i] = doc
	}

	result, err := r.traceCollection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to insert spans: %v", err)
	}

	log.Printf("Successfully inserted %d spans", len(result.InsertedIDs))
	return nil
}

func convertSpans(spans []models.Span) []bson.M {
	result := make([]bson.M, len(spans))
	for i, span := range spans {
		result[i] = bson.M{
			"id":                span.ID,
			"traceId":           span.TraceID,
			"spanId":            span.SpanID,
			"parentSpanId":      span.ParentSpanID,
			"name":              span.Name,
			"kind":              span.Kind,
			"startTimeUnixNano": span.StartTimeUnixNano,
			"endTimeUnixNano":   span.EndTimeUnixNano,
			"attributes":        convertAttributes(span.Attributes),
			"events":            convertEvents(span.Events),
			"links":             convertLinks(span.Links),
			"status":            convertStatus(span.Status),
		}
	}
	return result
}

func convertAttributes(attrs map[string]interface{}) bson.M {
	result := bson.M{}
	for k, v := range attrs {
		safeKey := strings.ReplaceAll(k, ".", "_")
		result[safeKey] = v
	}
	return result
}

func convertEvents(events []models.SpanEvent) []bson.M {
	result := make([]bson.M, len(events))
	for i, event := range events {
		result[i] = bson.M{
			"timeUnixNano": event.TimeUnixNano,
			"name":         event.Name,
			"attributes":   convertAttributes(event.Attributes),
		}
	}
	return result
}

func convertLinks(links []models.SpanLink) []bson.M {
	result := make([]bson.M, len(links))
	for i, link := range links {
		result[i] = bson.M{
			"traceId":    link.TraceID,
			"spanId":     link.SpanID,
			"attributes": convertAttributes(link.Attributes),
		}
	}
	return result
}

func convertStatus(status models.SpanStatus) bson.M {
	return bson.M{
		"code":        status.Code,
		"description": status.Description,
	}
}

func (r *MongoRepository) Close(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}
