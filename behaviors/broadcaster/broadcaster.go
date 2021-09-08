// Tideland Go Cells - Behaviors - Broadcaster
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package broadcaster // import "tideland.dev/go/cells/behaviors/broadcaster"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/cells/mesh"
)

//--------------------
// BEHAVIOR
//--------------------

// Behavior broadcasts all received events without change to all subscribers.
type Behavior struct{}

var _ mesh.Behavior = (*Behavior)(nil)

// New creates a broadcaster behavior.
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
			if err := out.EmitEvent(evt); err != nil {
				return err
			}
		}
	}
}

// EOF
