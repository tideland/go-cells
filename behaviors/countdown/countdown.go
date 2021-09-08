// Tideland Go Cells - Behaviors - Countdown
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package countdown // import "tideland.dev/go/cells/behaviors/countdown"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/cells/mesh"
)

//--------------------
// HELPER
//--------------------

// ZeroFunc is called when the countdown reaches zero. The collected
// events are passed, the returned event will be emitted.
type ZeroFunc func(r mesh.EventSinkReader) (*mesh.Event, error)

//--------------------
// BEHAVIOR
//--------------------

// Behavior collects a number of events. When this number is reached
// a zero function with access to these events will be called. The event
// returned by this function will be emitted and the counter reset.
type Behavior struct {
	t      int
	zeroer ZeroFunc
}

var _ mesh.Behavior = (*Behavior)(nil)

// New creates a countdown behavior based on the given t value
// and a zeroer function.
func New(t int, zeroer ZeroFunc) *Behavior {
	return &Behavior{
		t:      t,
		zeroer: zeroer,
	}
}

// Go implements the mesh.Behavior interface.
func (b *Behavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	sink := mesh.NewEventSink(b.t)
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			l := sink.Push(evt)
			if l == b.t {
				outEvt, err := b.zeroer(sink)
				if err != nil {
					return err
				}
				out.EmitEvent(outEvt)
				sink = mesh.NewEventSink(b.t)
			}
		}
	}
}

// EOF
