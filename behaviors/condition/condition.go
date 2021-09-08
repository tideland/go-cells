// Tideland Go Cells - Behaviors - Condition
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package condition // import "tideland.dev/go/cells/behaviors/condition"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/cells/mesh"
)

//--------------------
// HELPER
//--------------------

// ConditionTesterFunc checks if an event matches a wanted state.
type ConditionTesterFunc func(evt *mesh.Event) bool

// ConditionProcessorFunc handles the matching event.
type ConditionProcessorFunc func(cell mesh.Cell, evt *mesh.Event, out mesh.Emitter) error

//--------------------
// BEHAVIOR
//--------------------

// Behavior checks if an incoming event fillfills a given condition. This
// condition is defined by a given condition tester function. If that
// function returns true a given process function is called.
type Behavior struct {
	test    ConditionTesterFunc
	process ConditionProcessorFunc
}

var _ mesh.Behavior = (*Behavior)(nil)

// New creates a condition behavior instance with the given tester and
// process functions.
func New(tester ConditionTesterFunc, processor ConditionProcessorFunc) *Behavior {
	return &Behavior{
		test:    tester,
		process: processor,
	}
}

// Go implements the mesh.Behavior interface.
func (b *Behavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			if b.test(evt) {
				b.process(cell, evt, out)
			}
		}
	}
}

// EOF
