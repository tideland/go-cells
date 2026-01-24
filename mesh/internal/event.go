// Tideland Go Cells - Mesh - Internal
//
// Copyright (C) 2010-2026 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package internal // import "tideland.dev/go/cells/mesh/internal"

//--------------------
// IMPORTS
//--------------------

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

//--------------------
// METADATA
//--------------------

// Metadata contains additional information about an event for tracing,
// correlation, and custom data.
type Metadata struct {
	ID          string         // Unique event ID
	TraceID     string         // Trace ID for distributed tracing
	Correlation string         // Correlation ID for related events
	Custom      map[string]any // Custom metadata fields
}

// MetadataOption is a functional option for configuring event metadata.
type MetadataOption func(*Metadata)

// WithID sets a custom event ID.
func WithID(id string) MetadataOption {
	return func(m *Metadata) {
		m.ID = id
	}
}

// WithTraceID sets the trace ID for distributed tracing.
func WithTraceID(traceID string) MetadataOption {
	return func(m *Metadata) {
		m.TraceID = traceID
	}
}

// WithCorrelation sets the correlation ID for related events.
func WithCorrelation(correlation string) MetadataOption {
	return func(m *Metadata) {
		m.Correlation = correlation
	}
}

// WithCustom adds custom metadata fields.
func WithCustom(key string, value any) MetadataOption {
	return func(m *Metadata) {
		if m.Custom == nil {
			m.Custom = make(map[string]any)
		}
		m.Custom[key] = value
	}
}

// generateID generates a random unique ID for an event.
func generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp-based ID if random fails
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

//--------------------
// EVENT
//--------------------

// Event transports a topic and a payload a cell can process. The
// payload is anything marshalled into JSON and will be unmarshalled
// when a receiving cell accesses it.
type Event struct {
	timestamp time.Time
	emitters  []string
	topic     string
	payload   json.RawMessage
	metadata  Metadata
}

// NewEvent creates a new Event based on a topic. The payloads are optional.
// This is the legacy function maintained for backward compatibility.
func NewEvent(topic string, payloads ...any) (*Event, error) {
	if topic == "" {
		return nil, fmt.Errorf("event needs topic")
	}
	evt := &Event{
		timestamp: time.Now().UTC(),
		topic:     topic,
		metadata: Metadata{
			ID: generateID(),
		},
	}
	// Check if the only value is a payload.
	switch len(payloads) {
	case 0:
		return evt, nil
	case 1:
		bs, err := json.Marshal(payloads[0])
		if err != nil {
			return evt, fmt.Errorf("cannot marshal payload: %v", err)
		}
		evt.payload = bs
	default:
		bs, err := json.Marshal(payloads)
		if err != nil {
			return evt, fmt.Errorf("cannot marshal payload: %v", err)
		}
		evt.payload = bs
	}
	return evt, nil
}

// NewEventWithMetadata creates a new Event with metadata options.
func NewEventWithMetadata(topic string, payload any, opts ...MetadataOption) (*Event, error) {
	if topic == "" {
		return nil, fmt.Errorf("event needs topic")
	}
	metadata := Metadata{
		ID: generateID(),
	}
	for _, opt := range opts {
		opt(&metadata)
	}
	evt := &Event{
		timestamp: time.Now().UTC(),
		topic:     topic,
		metadata:  metadata,
	}
	if payload != nil {
		bs, err := json.Marshal(payload)
		if err != nil {
			return evt, fmt.Errorf("cannot marshal payload: %v", err)
		}
		evt.payload = bs
	}
	return evt, nil
}

// NewEventTyped creates a type-safe Event with a typed payload and metadata options.
// This is the preferred method for creating events with type safety.
func NewEventTyped[T any](topic string, payload T, opts ...MetadataOption) (*Event, error) {
	if topic == "" {
		return nil, fmt.Errorf("event needs topic")
	}
	metadata := Metadata{
		ID: generateID(),
	}
	for _, opt := range opts {
		opt(&metadata)
	}
	bs, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal payload: %v", err)
	}
	evt := &Event{
		timestamp: time.Now().UTC(),
		topic:     topic,
		payload:   bs,
		metadata:  metadata,
	}
	return evt, nil
}

// Timestamp returns the event timestamp.
func (evt Event) Timestamp() time.Time {
	return evt.timestamp
}

// Emitters returns an emitters path aloowing to see
// where an event has been emitted or simply re-emitted.
// The path layouts are
//
//     / is emitted via the mesh,
//     /foo is emitted by mesh and re-emitted by foo,
//     foo is emitted by foo,
//     foo/bar is emitted by foo and re-emitted by bar.
//
// So also longer paths like /foo/bar/baz are possible.
func (evt Event) Emitters() string {
	if len(evt.emitters) == 1 {
		return evt.emitters[0]
	}
	return evt.emitters[0] + strings.Join(evt.emitters[1:], "/")
}

// Topic returns the event topic.
func (evt Event) Topic() string {
	return evt.topic
}

// HasPayload checks if the event contains a payload.
func (evt Event) HasPayload() bool {
	return evt.payload != nil
}

// Payload unmarshals the payload of the event.
func (evt Event) Payload(payload any) error {
	if evt.payload == nil {
		return fmt.Errorf("Event contains no payload")
	}
	err := json.Unmarshal(evt.payload, payload)
	if err != nil {
		return fmt.Errorf("cannont unmarshal payload: %v", err)
	}
	return nil
}

// PayloadAs returns the typed payload of the event with compile-time type safety.
// This is the preferred method for accessing typed payloads.
func PayloadAs[T any](evt *Event) (T, error) {
	var zero T
	if evt.payload == nil {
		return zero, fmt.Errorf("Event contains no payload")
	}
	var result T
	err := json.Unmarshal(evt.payload, &result)
	if err != nil {
		return zero, fmt.Errorf("cannot unmarshal payload: %v", err)
	}
	return result, nil
}

// Metadata returns the event metadata.
func (evt Event) Metadata() Metadata {
	return evt.metadata
}

// String implements fmt.Stringer.
func (evt Event) String() string {
	return fmt.Sprintf(
		"Event{ID:%s Timestamp:%s Emitters:%v Topic:%v Payload:%v TraceID:%s Correlation:%s}",
		evt.metadata.ID,
		evt.timestamp.Format(time.RFC3339Nano),
		evt.emitters,
		evt.topic,
		string(evt.payload),
		evt.metadata.TraceID,
		evt.metadata.Correlation,
	)
}

// MarshalJSON implements the custom JSON marshaling of the event.
func (evt Event) MarshalJSON() ([]byte, error) {
	tmp := struct {
		Timestamp time.Time       `json:"timestamp"`
		Emitters  []string        `json:"emitters,omitempty"`
		Topic     string          `json:"topic"`
		Payload   json.RawMessage `json:"payload,omitempty"`
		Metadata  Metadata        `json:"metadata"`
	}{
		Timestamp: evt.timestamp,
		Emitters:  evt.emitters,
		Topic:     evt.topic,
		Payload:   evt.payload,
		Metadata:  evt.metadata,
	}
	return json.Marshal(tmp)
}

// UnmarshalJSON implements the custom JSON unmarshaling of the event.
func (evt *Event) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Timestamp time.Time       `json:"timestamp"`
		Emitters  []string        `json:"emitters,omitempty"`
		Topic     string          `json:"topic"`
		Payload   json.RawMessage `json:"payload,omitempty"`
		Metadata  Metadata        `json:"metadata"`
	}{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	evt.timestamp = tmp.Timestamp
	evt.emitters = tmp.Emitters
	evt.topic = tmp.Topic
	evt.payload = tmp.Payload
	evt.metadata = tmp.Metadata
	return nil
}

// initEmitters sets the emitters to the mesh value.
func (evt *Event) initEmitters() {
	evt.emitters = []string{"/"}
}

// appendEmitter is used by the different emitters to signal their a
// sender or passer of an event.
func (evt *Event) appendEmitter(name string) {
	evt.emitters = append(evt.emitters, name)
}

// EOF
