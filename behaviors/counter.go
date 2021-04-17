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
// COUNTER BEHAVIOR
//--------------------

// CounterEvaluationFunc analyzes the passed event and returns, which counters
// shall be incremented.
type CounterEvaluationFunc func(evt *mesh.Event) ([]string, error)

// counterBehavior counts events based on the counter function.
type counterBehavior struct {
	eval   CounterEvaluationFunc
	values map[string]int
}

// NewCounterBehavior creates a counter behavior based on the passed
// function. That function has to evaluate incomming events and to return
// the names of counters to increase. All counters start at 0. An
// event with the topic "counter-status" emits all current counters,
// an event with the topic "counter-reset!" clears them.
func NewCounterBehavior(eval CounterEvaluationFunc) mesh.Behavior {
	return &counterBehavior{
		eval:   eval,
		values: make(map[string]int),
	}
}

// Go evaluates the incoming events and increases according variables.
func (b *counterBehavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			switch evt.Topic() {
			case TopicCounterReset:
				b.values = make(map[string]int)
			case TopicCounterStatus:
				if err := out.Emit(TopicCounterValues, b.values); err != nil {
					return err
				}
			default:
				incrs, err := b.eval(evt)
				if err != nil {
					return err
				}
				for _, incr := range incrs {
					b.values[incr]++
				}
			}
		}
	}
}

// EOF
