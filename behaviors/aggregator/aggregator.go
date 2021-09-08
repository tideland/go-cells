// Tideland Go Cells - Behaviors - Aggregator
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package aggregator // import "tideland.dev/go/cells/behaviors/aggregator"

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
	TopicAggregate     = "aggregate!"
	TopicAggregateDone = "aggregate-done"
	TopicReset         = "reset!"
	TopicResetDone     = "reset-done"
)

//--------------------
// HELPER
//--------------------

// AggregatorFunc is a function receiving the current status payload
// and event and returns the next status payload.
type AggregatorFunc func(status interface{}, evt *mesh.Event) (interface{}, error)

//--------------------
// BEHAVIOR
//--------------------

// Behavior provides a behavior which aggregates the stream of events with
// a given function. A received "reset!" topic resets the status.
type Behavior struct {
	initialize func() interface{}
	status     interface{}
	aggregate  AggregatorFunc
}

var _ mesh.Behavior = (*Behavior)(nil)

// New creates an instance of the aggregator behavior with the given
// aggregator function.The initializer function creates the first value
// before aggregating.
func New(initializer func() interface{}, aggregator AggregatorFunc) *Behavior {
	b := &Behavior{
		initialize: initializer,
		status:     nil,
		aggregate:  aggregator,
	}
	if b.initialize != nil {
		b.status = b.initialize()
	}
	return b
}

// Go implements the mesh.Behavior interface.
func (b *Behavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			switch evt.Topic() {
			case TopicAggregate:
				if err := out.Emit(TopicAggregateDone, b.status); err != nil {
					return err
				}
			case TopicReset:
				b.status = nil
				if b.initialize != nil {
					b.status = b.initialize()
				}
				if err := out.Emit(TopicResetDone, b.status); err != nil {
					return err
				}
			default:
				status, err := b.aggregate(b.status, evt)
				if err != nil {
					return err
				}
				b.status = status
			}
		}
	}
}

// EOF
