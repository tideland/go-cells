// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package filter

import (
	"tideland.dev/go/cells/mesh"
)

const (
	TopicEvaluate       = "evaluate!"
	TopicEvaluationDone = "evaluation-done"
	TopicReset          = "reset!"
	TopicResetDone      = "reset-done"
)

// FilterFunc defines how events are filtered for including or excluding.
type FilterFunc func(event *mesh.Event) (bool, error)

// mode describes if the filter works including or excluding.
type mode int

// Flags for the filter mode.
const (
	includingMode mode = iota
	excludingMode
)

// Behavior provides a behavior allowing to filter the stream of incomming
// events in a user defined was. The way instantiating it defines if the
// filter function decides if events are included or excluded.
type Behavior struct {
	filter FilterFunc
	mode   mode
}

var _ mesh.Behavior = (*Behavior)(nil)

// NewIncluding creates a new instance of the filter including those
// events where the given filter function returns true.
func NewIncluding(filter FilterFunc) *Behavior {
	return &Behavior{
		filter: filter,
		mode:   includingMode,
	}
}

// NewExcluding creates a new instance of the filter including those
// events where the given filter function returns true.
func NewExcluding(filter FilterFunc) *Behavior {
	return &Behavior{
		filter: filter,
		mode:   excludingMode,
	}
}

// Go implements the mesh.Behavior interface.
func (b *Behavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			ok, err := b.filter(evt)
			if err != nil {
				return err
			}
			switch b.mode {
			case includingMode:
				if ok {
					out.EmitEvent(evt)
				}
			case excludingMode:
				if !ok {
					out.EmitEvent(evt)
				}
			}
		}
	}
}
