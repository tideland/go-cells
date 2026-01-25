// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package mapper

import (
	"tideland.dev/go/cells/mesh"
)

// MapperFunc defines how incoming events are to map into outcomming event.
type MapperFunc func(evt *mesh.Event) (*mesh.Event, error)

// Behavior provides a behavior allowing to map incoming events into outgohing
// events, e.g. with different topics based on an analysis, or with a changed
// payload. The mapping is defined via the mapper function which also may return
// nil. In this case no event is emitted.
type Behavior struct {
	mapper MapperFunc
}

// New creates a new instance of the mapper behavior.
func New(mapper MapperFunc) *Behavior {
	return &Behavior{
		mapper: mapper,
	}
}

// Go implements the mesh.Behavior interface.
func (b *Behavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			mapped, err := b.mapper(evt)
			if err != nil {
				return err
			}
			if mapped != nil {
				out.EmitEvent(mapped)
			}
		}
	}
}
