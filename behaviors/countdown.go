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
// COUNTDOWN BEHAVIOR
//--------------------

// Zeroer is called when the countdown reaches zero. The collected
// events are passed, the returned event will be emitted.
type Zeroer func(r mesh.EventSinkReader) (*mesh.Event, error)

// CountdownBehavior collects a number of events. When this number is reached
// a zero function with access to these eventsn will be called. Its returned
// event will be emitted.
type CountdownBehavior struct {
	t      int
	zeroer Zeroer
}

// NewCountdownBehavior creates a countdown behavior based on the passed
// t value and zeroer function.
func NewCountdownBehavior(t int, zeroer Zeroer) *CountdownBehavior {
	return &CountdownBehavior{
		t:      t,
		zeroer: zeroer,
	}
}

// Go counts and collects received events for processing them en bloc.
func (b *CountdownBehavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
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
			}
		}
	}
}

// EOF
