// Tideland Go Cells - Behaviors - Callback
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package callback // import "tideland.dev/go/cells/behaviors/callback"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/cells/mesh"
)

//--------------------
// HELPER
//--------------------

// CallbackFunc is a function called by the behavior when it receives an event.
type CallbackFunc func(evt *mesh.Event, out mesh.Emitter) error

//--------------------
// BEHAVIOR
//--------------------

// Behavior implements a behavior calling a muber of functions for
// each event.
type Behavior struct {
	callbacks []CallbackFunc
}

var _ mesh.Behavior = &Behavior{}

// New creates an instance using the given callback functions.
func New(callbacks ...CallbackFunc) *Behavior {
	return &Behavior{
		callbacks: callbacks,
	}
}

// Go implements the mesh.Behavior interface.
func (b *Behavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			for _, callback := range b.callbacks {
				if err := callback(evt, out); err != nil {
					return err
				}
			}
		}
	}
}

// EOF
