// Tideland Go Cells - Behaviors - Filter
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package filter // import "tideland.dev/go/cells/behaviors/filter"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TOPICS
//--------------------

const (
	TopicEvaluate       = "evaluate!"
	TopicEvaluationDone = "evaluation-done"
	TopicReset          = "reset!"
	TopicResetDone      = "reset-done"
)

//--------------------
// HELPER
//--------------------

// FilterFunc defines how events are filtered for including or excluding.
type FilterFunc func(event *mesh.Event) (bool, error)

// mode describes if the filter works including or excluding.
type mode int

// Flags for the filter mode.
const (
	includingMode mode = iota
	excludingMode
)

//--------------------
// BEHAVIOR
//--------------------

// Behavior
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

// EOF
