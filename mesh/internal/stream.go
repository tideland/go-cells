// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package internal

import (
	"errors"
	"time"
)

// Stream manages the flow of events between emitter and receiver.
type streamImpl struct {
	eventc chan *Event
}

// NewStream creates a stream instance.
func NewStream() *streamImpl {
	return &streamImpl{
		eventc: make(chan *Event),
	}
}

// Pull reads an event out of the stream.
func (str *streamImpl) Pull() <-chan *Event {
	return str.eventc
}

// Emit creates a new event and emits it.
func (str *streamImpl) Emit(topic string, payloads ...any) error {
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
func (str *streamImpl) EmitEvent(evt *Event) error {
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
