// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package internal

import (
	"context"
)

// Mesh describes the interface to a mesh of a cell from the
// perspective of a behavior.
type Mesh interface {
	// Go starts a cell using the given behavior.
	Go(name string, b Behavior) error

	// Subscribe subscribes the cell with receptor name to the cell
	// with emitter name.
	Subscribe(emitterName, receptorName string) error

	// Unsubscribe unsubscribes the cell with receptor name from the cell
	// with emitter name.
	Unsubscribe(emitterName, receptorName string) error

	// Emit creates an event and raises it to the named cell.
	Emit(name, topic string, payloads ...any) error

	// EmitEvent raises an event to the named cell.
	EmitEvent(name string, evt *Event) error

	// Emitter returns a static emitter for the named cell.
	Emitter(name string) (Emitter, error)
}

// Cell describes the interface to a cell from the perspective
// of a behavior.
type Cell interface {
	// Context returns the context of mesh and cell.
	Context() context.Context

	// Name returns the name of the deployed cell running the
	// behavior.
	Name() string

	// Mesh returns the mesh of the cell.
	Mesh() Mesh
}

// Behavior describes what cell implementations must understand.
type Behavior interface {
	// Go will be started as wrapped goroutine.
	Go(cell Cell, in Receptor, out Emitter) error
}

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

// EventSinkDoFunc is used when looking over the collected events.
type EventSinkDoFunc func(i int, evt *Event) error

// EventSinkChanger can be used to write events into a sink or change
// it by reading and deleting.
type EventSinkChanger interface {
	// Push adds an event to the end of the sink and returns the new size.
	Push(evt *Event) int

	// Pop retrieves and removes the last event from the sink
	// and also returns the new length.
	Pop() (*Event, int)

	// Unshift adds an event to the begin of the sink and returns the new size.
	Unshift(evt *Event) int

	// Shift returns and removes the first event of the sink
	// and also returns the new length.
	Shift() (*Event, int)

	// Clear removes all collected events.
	Clear()
}

// EventSinkReader can be used to read the events in a sink.
type EventSinkReader interface {
	// Len returns the number of stored events.
	Len() int

	// First returns the first of the collected events.
	First() (*Event, bool)

	// Last returns the last of the collected event datas.
	Last() (*Event, bool)

	// Peek returns an event at a given index and true if it
	// exists, otherwise nil and false.
	Peek(index int) (*Event, bool)

	// Do iterates over all collected events.
	Do(do EventSinkDoFunc) error

	// String returns a string representation.
	String() string
}

// EventSink combines changer and reader. It stores a number of ordered events by
// adding them at the end. To be used in behaviors for collecting sets of events
// and operate on them.
type EventSink interface {
	EventSinkChanger
	EventSinkReader
}
