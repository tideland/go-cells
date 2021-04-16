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
// CONDITION BEHAVIOR
//--------------------

// ConditionTesterFunc checks if an event matches a wanted state.
type ConditionTesterFunc func(evt mesh.Event) bool

// ConditionProcessorFunc handles the matching event.
type ConditionProcessorFunc func(cell mesh.Cell, evt mesh.Event, out mesh.Emitter) error

// conditionBehavior implements the condition behavior.
type conditionBehavior struct {
	test    ConditionTesterFunc
	process ConditionProcessorFunc
}

// NewConditionBehavior creates a behavior testing of a cell
// fullfills a given condition. If the test returns true the
// processor is called.
func NewConditionBehavior(tester ConditionTesterFunc, processor ConditionProcessorFunc) mesh.Behavior {
	return &conditionBehavior{
		test:    tester,
		process: processor,
	}
}

// Go checks the condition and calls the process in case of a
// positive test.
func (b *conditionBehavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
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
