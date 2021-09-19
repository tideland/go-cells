// Tideland Go Cells - Behaviors - Rate Evaluator
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package rateevaluator // import "tideland.dev/go/cells/behaviors/rateevaluator"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/cells/mesh"
)

//--------------------
// CONSTANTS
//--------------------

const (
	TopicDummy = "dummy"
)

//--------------------
// HELPER
//--------------------

//--------------------
// BEHAVIOR
//--------------------

// Behavior provides a behavior ...
type Behavior struct {
}

// New creates a new instance of the rate evaluator behavior.
func New() *Behavior {
	return &Behavior{}
}

// Go implements the mesh.Behavior interface.
func (b *Behavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
		}
	}
}

// EOF
