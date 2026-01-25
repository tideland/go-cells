// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package broadcaster

import (
	"tideland.dev/go/cells/mesh"
)

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
