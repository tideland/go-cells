// Tideland Go Cells - Behaviors - Combo
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package combo // import "tideland.dev/go/cells/behaviors/combo"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"

	"tideland.dev/go/cells/mesh"
)

//--------------------
// TOPICS
//--------------------

const (
	TopicCriterionDone = "criterion-done"
	TopicReset         = "reset!"
	TopicResetDone     = "reset-done"
)

//--------------------
// HELPER
//--------------------

// ComboCriterionFunc is used by the combo behavior. It has to return
// CriterionDone when a combination is complete, CriterionKeep when it
// is so far okay but not yet complete, CriterionDropFirst when the first
// event shall be dropped, CriterionDropLast when the last event shall
// be dropped, and CriterionClear when the collected events have to be
// cleared for starting over. In case of CriterionDone it additionally
// has to return a payload which will be emitted.
type ComboCriterionFunc func(r mesh.EventSinkReader) (CriterionMatch, interface{}, error)

// CriterionMatch allows a combo criterion func to signal its
// analysis rersult.
type CriterionMatch int

// Criterion matches.
const (
	CriterionError CriterionMatch = iota
	CriterionDone
	CriterionKeep
	CriterionDropFirst
	CriterionDropLast
)

//--------------------
// BEHAVIOR
//--------------------

// Behavior checks the event stream for a combination of events defined by
// a criterion function. In case of a match an according event is emitted.
type Behavior struct {
	matches ComboCriterionFunc
	sink    mesh.EventSink
}

var _ mesh.Behavior = (*Behavior)(nil)

// New creates an instance of the combo behavior using the given criterion
// function.
func New(matcher ComboCriterionFunc) *Behavior {
	return &Behavior{
		matches: matcher,
		sink:    mesh.NewEventSink(0),
	}
}

// Go implements the mesh.Behavior interface.
func (b *Behavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			switch evt.Topic() {
			case TopicReset:
				b.sink.Clear()
				out.Emit(TopicResetDone)
			default:
				b.sink.Push(evt)
				matches, data, err := b.matches(b.sink)
				if err != nil {
					return err
				}
				switch matches {
				case CriterionDone:
					out.Emit(TopicCriterionDone, data)
					b.sink.Clear()
				case CriterionKeep:
				case CriterionDropFirst:
					b.sink.Shift()
				case CriterionDropLast:
					b.sink.Pop()
				default:
					return fmt.Errorf("invalid criterion matcher result: %v", matches)
				}
			}
		}
	}
}

// EOF
