// Tideland Go Cells - Behaviors - Collector
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package collector // import "tideland.dev/go/cells/behaviors/collector"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TOPICS
//--------------------

const (
	TopicProcess        = "process!"
	TopicProcessingDone = "processing-done"
	TopicReset          = "reset!"
	TopicResetDone      = "reset-done"
)

//--------------------
// HELPER
//--------------------

// CollectionProcessorFunc is used to process collected events.
type CollectionProcessorFunc func(r mesh.EventSinkReader) (*mesh.Event, error)

//--------------------
// BEHAVIOR
//--------------------

// Behavior collects a number events for processing on demand.
type Behavior struct {
	max     int
	sink    mesh.EventSink
	process CollectionProcessorFunc
}

var _ mesh.Behavior = &Behavior{}

// New collects geven maximum number of events. If the number gets too large the first
// one will be deleted. After "process!" topic it processes it and emits the result as
// event. After "reset!" topic the collection is dropped to zero.
func New(max int, process CollectionProcessorFunc) *Behavior {
	return &Behavior{
		max:     max,
		sink:    mesh.NewEventSink(max),
		process: process,
	}
}

// Go implements the mesh.Behavior interface.
func (b *Behavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return cell.Context().Err()
		case evt := <-in.Pull():
			switch evt.Topic() {
			case TopicReset:
				b.sink.Clear()
				out.Emit(TopicResetDone)
			case TopicProcess:
				pevt, err := b.process(b.sink)
				if err != nil {
					return err
				}
				out.EmitEvent(pevt)
				b.sink.Clear()
			default:
				b.sink.Push(evt)
			}
		}
	}
}

// EOF
