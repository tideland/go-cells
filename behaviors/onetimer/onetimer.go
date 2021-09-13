// Tideland Go Cells - Behaviors - One-Timer
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package onetimer // import "tideland.dev/go/cells/behaviors/onetimer"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/cells/mesh"
)

//--------------------
// HELPER
//--------------------

// OneTimerFunc describes the function called after the first event.
type OneTimerFunc func(evt *mesh.Event, out mesh.Emitter) error

//--------------------
// BEHAVIOR
//--------------------

// Behavior implements a behavior calling the one-timer function the
// first time it receives any event. This user-defined function can analyze
// the event and spawn new events. Afterwards it will process no received
// event anymore.
type Behavior struct {
	oneTime OneTimerFunc
}

// New creates a one-time behavior using the given function.
func New(oneTime OneTimerFunc) *Behavior {
	return &Behavior{
		oneTime: oneTime,
	}
}

// Go implements the mesh.Behavior interface.
func (b *Behavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			if b.oneTime != nil {
				if err := b.oneTime(evt, out); err != nil {
					return err
				}
				b.oneTime = nil
			}
		}
	}
}

// EOF
