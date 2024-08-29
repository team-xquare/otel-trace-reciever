package telemetry

import (
	"encoding/hex"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
	"math"
	"otel-trace-reciever/internal/models"

	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

func ConvertResourceSpansToTraces(resourceSpans *tracepb.ResourceSpans) []*models.Trace {
	traceMap := make(map[string]*models.Trace)
	serviceName := getServiceName(resourceSpans.Resource)

	for _, scopeSpans := range resourceSpans.ScopeSpans {
		for _, span := range scopeSpans.Spans {
			convertedSpan := convertSpan(span)
			traceID := convertedSpan.TraceID

			if trace, exists := traceMap[traceID]; exists {
				trace.Spans = append(trace.Spans, convertedSpan)
				updateTraceDuration(trace, convertedSpan)
			} else {
				traceMap[traceID] = &models.Trace{
					TraceID:      traceID,
					Spans:        []models.Span{convertedSpan},
					ServiceName:  serviceName,
					DateNano:     convertedSpan.StartTimeUnixNano,
					DurationNano: convertedSpan.EndTimeUnixNano - convertedSpan.StartTimeUnixNano,
				}
			}
		}
	}

	var traces []*models.Trace
	for _, trace := range traceMap {
		traces = append(traces, trace)
	}

	return traces
}

func updateTraceDuration(trace *models.Trace, span models.Span) {
	if span.StartTimeUnixNano < trace.DateNano {
		trace.DateNano = span.StartTimeUnixNano
	}
	if span.EndTimeUnixNano-trace.DateNano > trace.DurationNano {
		trace.DurationNano = span.EndTimeUnixNano - trace.DateNano
	}
}

func convertSpan(pbSpan *tracepb.Span) models.Span {
	return models.Span{
		ID:                hex.EncodeToString(pbSpan.SpanId),
		TraceID:           hex.EncodeToString(pbSpan.TraceId),
		SpanID:            hex.EncodeToString(pbSpan.SpanId),
		ParentSpanID:      convertParentSpanID(pbSpan.ParentSpanId),
		Name:              pbSpan.Name,
		Kind:              int(pbSpan.Kind),
		StartTimeUnixNano: safeUint64ToInt64(pbSpan.StartTimeUnixNano),
		EndTimeUnixNano:   safeUint64ToInt64(pbSpan.EndTimeUnixNano),
		Attributes:        convertAttributesPb(pbSpan.Attributes),
		Events:            convertEvents(pbSpan.Events),
		Links:             convertLinks(pbSpan.Links),
		Status:            convertStatus(pbSpan.Status),
	}
}

func safeUint64ToInt64(value uint64) int64 {
	if value > math.MaxInt64 {
		return math.MaxInt64
	}
	return int64(value)
}

func convertParentSpanID(parentSpanID []byte) *string {
	if len(parentSpanID) == 0 {
		return nil
	}
	s := hex.EncodeToString(parentSpanID)
	return &s
}

func convertAttributesPb(attrs []*commonpb.KeyValue) map[string]interface{} {
	attributes := make(map[string]interface{})
	for _, attr := range attrs {
		switch v := attr.Value.Value.(type) {
		case *commonpb.AnyValue_StringValue:
			attributes[attr.Key] = v.StringValue
		case *commonpb.AnyValue_BoolValue:
			attributes[attr.Key] = v.BoolValue
		case *commonpb.AnyValue_IntValue:
			attributes[attr.Key] = v.IntValue
		case *commonpb.AnyValue_DoubleValue:
			attributes[attr.Key] = v.DoubleValue
		case *commonpb.AnyValue_ArrayValue:
			attributes[attr.Key] = convertArrayValue(v.ArrayValue)
		case *commonpb.AnyValue_KvlistValue:
			attributes[attr.Key] = convertKeyValueList(v.KvlistValue)
		}
	}
	return attributes
}

func convertArrayValue(arrayValue *commonpb.ArrayValue) []interface{} {
	result := make([]interface{}, len(arrayValue.Values))
	for i, value := range arrayValue.Values {
		result[i] = convertAnyValue(value)
	}
	return result
}

func convertKeyValueList(kvList *commonpb.KeyValueList) map[string]interface{} {
	result := make(map[string]interface{})
	for _, kv := range kvList.Values {
		result[kv.Key] = convertAnyValue(kv.Value)
	}
	return result
}

func convertAnyValue(value *commonpb.AnyValue) interface{} {
	switch v := value.Value.(type) {
	case *commonpb.AnyValue_StringValue:
		return v.StringValue
	case *commonpb.AnyValue_BoolValue:
		return v.BoolValue
	case *commonpb.AnyValue_IntValue:
		return v.IntValue
	case *commonpb.AnyValue_DoubleValue:
		return v.DoubleValue
	case *commonpb.AnyValue_ArrayValue:
		return convertArrayValue(v.ArrayValue)
	case *commonpb.AnyValue_KvlistValue:
		return convertKeyValueList(v.KvlistValue)
	default:
		return nil
	}
}

func convertEvents(events []*tracepb.Span_Event) []models.SpanEvent {
	var spanEvents []models.SpanEvent
	for _, event := range events {
		spanEvents = append(spanEvents, models.SpanEvent{
			TimeUnixNano: safeUint64ToInt64(event.TimeUnixNano),
			Name:         event.Name,
			Attributes:   convertAttributesPb(event.Attributes),
		})
	}
	return spanEvents
}

func convertLinks(links []*tracepb.Span_Link) []models.SpanLink {
	var spanLinks []models.SpanLink
	for _, link := range links {
		spanLinks = append(spanLinks, models.SpanLink{
			TraceID:    hex.EncodeToString(link.TraceId),
			SpanID:     hex.EncodeToString(link.SpanId),
			Attributes: convertAttributesPb(link.Attributes),
		})
	}
	return spanLinks
}

func convertStatus(status *tracepb.Status) models.SpanStatus {
	if status == nil {
		return models.SpanStatus{}
	}
	return models.SpanStatus{
		Code:        int(status.Code),
		Description: status.Message,
	}
}

func getServiceName(resource *resourcepb.Resource) *string {
	for _, attr := range resource.Attributes {
		if attr.Key == "service.name" {
			value := attr.Value.GetStringValue()
			return &value
		}
	}
	return nil
}