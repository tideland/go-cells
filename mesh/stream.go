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
	"errors"
	"time"
)

//--------------------
// INTERFACES
//--------------------

// Receptor defines the interface to receive events.
type Receptor interface {
	// Pull reads an event out of the input stream.
	Pull() <-chan *Event
}

// Emitter defines the interface for emitting events to one
// or more cells.
type Emitter interface {
	// Emit creates a new event and appends it to the output stream.
	Emit(topic string, payloads ...any) error

	// EmitEvent appends the given event to the output stream.
	EmitEvent(evt *Event) error
}

//--------------------
// STREAM
//--------------------

// stream manages the flow of events between emitter and receiver.
type stream struct {
	eventc chan *Event
}

// newStream creates a stream instance.
func newStream() *stream {
	return &stream{
		eventc: make(chan *Event),
	}
}

// Pull reads an event out of the stream.
func (str *stream) Pull() <-chan *Event {
	return str.eventc
}

// Emit creates a new event and emits it.
func (str *stream) Emit(topic string, payloads ...any) error {
	evt, err := NewEvent(topic, payloads...)
	if err != nil {
		return err
	}
	return str.EmitEvent(evt)
}

// EmitEvent appends an event to the end of the stream. It retries to
// append it to the buffer in case that it's full. The time will
// increase. If waiting is longer than 5 seconds a timeout error will
// be returned.
func (str *stream) EmitEvent(evt *Event) error {
	total := 5 * time.Second
	wait := 50 * time.Millisecond
	waited := 0 * time.Millisecond
	for {
		select {
		case str.eventc <- evt:
			return nil
		default:
			time.Sleep(wait)
			waited += wait
			wait += 50 * time.Millisecond
			if waited > total {
				return errors.New("timeout")
			}
		}
	}
}

// EOF
