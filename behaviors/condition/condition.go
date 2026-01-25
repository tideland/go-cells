// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package condition

import (
	"tideland.dev/go/cells/mesh"
)

// ConditionTesterFunc checks if an event matches a wanted state.
type ConditionTesterFunc func(evt *mesh.Event) bool

// ConditionProcessorFunc handles the matching event.
type ConditionProcessorFunc func(cell mesh.Cell, evt *mesh.Event, out mesh.Emitter) error

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
