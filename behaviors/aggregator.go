// Tideland Go Cells - Behaviors
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors // import "tideland.dev/go/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/cells/mesh"
)

//--------------------
// AGGREGATOR BEHAVIOR
//--------------------

// AggregatorFunc is a function receiving the current status payload
// and event and returns the next status payload.
type AggregatorFunc func(status interface{}, evt *mesh.Event) (interface{}, error)

// aggregatorBehavior implements the aggregator behavior.
type aggregatorBehavior struct {
	init      func() interface{}
	status    interface{}
	aggregate AggregatorFunc
}

// NewAggregatorBehavior creates a behavior aggregating the received events
// and emits events with the new aggregate. A "reset!" topic resets the
// aggregate to nil again.
func NewAggregatorBehavior(initializer func() interface{}, aggregator AggregatorFunc) mesh.Behavior {
	b := &aggregatorBehavior{
		init:      initializer,
		status:    nil,
		aggregate: aggregator,
	}
	if b.init != nil {
		b.status = b.init()
	}
	return b
}

// Go aggregates the event.
func (b *aggregatorBehavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			switch evt.Topic() {
			case TopicAggregate:
				// Request status.
				if err := out.Emit(TopicAggregated, b.status); err != nil {
					return err
				}
			case TopicReset:
				// Reset.
				b.status = nil
				if b.init != nil {
					b.status = b.init()
				}
				if err := out.Emit(TopicResetted, b.status); err != nil {
					return err
				}
			default:
				// Aggregate the event.
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
