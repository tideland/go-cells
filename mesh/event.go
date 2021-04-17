// Tideland Go Cells - Mesh
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh // import "tideland.dev/go/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

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
}

// NewEvent creates a new Event based on a topic. The payloads are optional.
func NewEvent(topic string, payloads ...interface{}) (*Event, error) {
	if topic == "" {
		return nil, fmt.Errorf("event needs topic")
	}
	evt := &Event{
		timestamp: time.Now().UTC(),
		topic:     topic,
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
func (evt Event) Payload(payload interface{}) error {
	if evt.payload == nil {
		return fmt.Errorf("Event contains no payload")
	}
	err := json.Unmarshal(evt.payload, payload)
	if err != nil {
		return fmt.Errorf("cannont unmarshal payload: %v", err)
	}
	return nil
}

// String implements fmt.Stringer.
func (evt Event) String() string {
	return fmt.Sprintf(
		"Event{Timestamp:%s Emitters:%v Topic:%v Payload:%v}",
		evt.timestamp.Format(time.RFC3339Nano),
		evt.emitters,
		evt.topic,
		string(evt.payload),
	)
}

// MarshalJSON implements the custom JSON marshaling of the event.
func (evt Event) MarshalJSON() ([]byte, error) {
	tmp := struct {
		Timestamp time.Time       `json:"timestamp"`
		Emitters  []string        `json:"emitters,omitempty"`
		Topic     string          `json:"topic"`
		Payload   json.RawMessage `json:"payload,omitempty"`
	}{
		Timestamp: evt.timestamp,
		Emitters:  evt.emitters,
		Topic:     evt.topic,
		Payload:   evt.payload,
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
	}{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	evt.timestamp = tmp.Timestamp
	evt.emitters = tmp.Emitters
	evt.topic = tmp.Topic
	evt.payload = tmp.Payload
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
