package models

type Trace struct {
	TraceID      string  `bson:"traceId"`
	Spans        []Span  `bson:"spans"`
	ServiceName  *string `bson:"serviceName,omitempty"`
	DateNano     int64   `bson:"dateNano"`
	DurationNano int64   `bson:"durationNano"`
}

type Span struct {
	ID                string                 `bson:"id"`
	TraceID           string                 `bson:"traceId"`
	SpanID            string                 `bson:"spanId"`
	ParentSpanID      *string                `bson:"parentSpanId,omitempty"`
	Name              string                 `bson:"name"`
	Kind              int                    `bson:"kind"`
	StartTimeUnixNano int64                  `bson:"startTimeUnixNano"`
	EndTimeUnixNano   int64                  `bson:"endTimeUnixNano"`
	Attributes        map[string]interface{} `bson:"attributes"`
	Events            []SpanEvent            `bson:"events"`
	Links             []SpanLink             `bson:"links"`
	Status            SpanStatus             `bson:"status"`
}

type SpanEvent struct {
	TimeUnixNano int64                  `bson:"timeUnixNano"`
	Name         string                 `bson:"name"`
	Attributes   map[string]interface{} `bson:"attributes"`
}

type SpanLink struct {
	TraceID    string                 `bson:"traceId"`
	SpanID     string                 `bson:"spanId"`
	Attributes map[string]interface{} `bson:"attributes"`
}

type SpanStatus struct {
	Code        int    `bson:"code"`
	Description string `bson:"description,omitempty"`
}
