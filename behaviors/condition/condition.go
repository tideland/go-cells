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

// Behavior check if an incoming event fillfills a given condition. If the test
// function returns true the process function is called.
type Behavior struct {
	test    ConditionTesterFunc
	process ConditionProcessorFunc
}

var _ mesh.Behavior = &Behavior{}

// New creates a behavior testing of a cell
// fullfills a given condition. If the test returns true the
// processor is called.
func New(tester ConditionTesterFunc, processor ConditionProcessorFunc) *Behavior {
	return &Behavior{
		test:    tester,
		process: processor,
	}
}

// Go checks the condition and calls the process in case of a
// positive test.
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
