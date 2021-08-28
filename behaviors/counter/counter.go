// Tideland Go Cells - Behaviors - Counter
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package counter // import "tideland.dev/go/cells/behaviors/counter"

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
	TopicCounters     = "counters!"
	TopicCountersDone = "counters-done!"
	TopicReset        = "reset!"
	TopicResetDone    = "reset-done"
)

//--------------------
// HELPER
//--------------------

// CounterEvaluationFunc analyzes the passed event and returns, which counters
// shall be incremented.
type CounterEvaluationFunc func(evt *mesh.Event) ([]string, error)

//--------------------
// BEHAVIOR
//--------------------

// Behavior evaluates incoming events with a given CounterEvaluationFunc. This
// function decides by returning a number of identifiers, which counter will
// be incremented. All counters can be reset with the topic "reset!" and the
// counters sent by "counters!".
type Behavior struct {
	eval    CounterEvaluationFunc
	counter map[string]int
}

var _ mesh.Behavior = &Behavior{}

// New instantiatas a counter behavior with the given evaluator.
func New(eval CounterEvaluationFunc) *Behavior {
	return &Behavior{
		eval:    eval,
		counter: make(map[string]int),
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
				b.counter = make(map[string]int)
				out.Emit(TopicResetDone)
			case TopicCounters:
				if err := out.Emit(TopicCountersDone, b.counter); err != nil {
					return err
				}
			default:
				incrs, err := b.eval(evt)
				if err != nil {
					return err
				}
				for _, incr := range incrs {
					b.counter[incr]++
				}
			}
		}
	}
}

// EOF
