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
	"time"

	"tideland.dev/go/cells/mesh"
)

//--------------------
// CONSTANTS
//--------------------

const (
	TopicRate  = "rate"
	TopicReset = "reset!"
)

//--------------------
// HELPER
//--------------------

// RaterFunc is used by the rate evaluator behavior and has to return true, if
// the given event matches a criterion for rate measuring.
type RaterFunc func(evt *mesh.Event) (bool, error)

// Rate describes the rate of events matching the given criterion. It
// contains the matching time, the duration from the last match to this
// one, and the highest, lowest, and avaerage duration between matches.
type Rate struct {
	Time             time.Time
	CountMatching    int
	CountNonMatching int
	Duration         time.Duration
	High             time.Duration
	Low              time.Duration
	Average          time.Duration
}

//--------------------
// BEHAVIOR
//--------------------

// Behavior provides a behavior evaluating event rates. Each time a rater
// func returns true for a received event the duration between this and the
// last one is measured and logged. Also a rate payload is emitted.
type Behavior struct {
	matches RaterFunc
	count   int
}

var _ mesh.Behavior = (*Behavior)(nil)

// New creates an event rate evaluating behavior.
func New(matches RaterFunc, count int) *Behavior {
	return &Behavior{
		matches: matches,
		count:   count,
	}
}

// Go implements the mesh.Behavior interface.
func (b *Behavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	last := time.Now()
	countMatching := 0
	countNonMatching := 0
	durations := []time.Duration{}
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			switch evt.Topic() {
			case TopicReset:
				last = time.Now()
				countMatching = 0
				countNonMatching = 0
				durations = []time.Duration{}
			default:
				ok, err := b.matches(evt)
				if err != nil {
					return err
				}
				if !ok {
					countNonMatching++
					continue
				}
				// Recalculate rate with matching event.
				countMatching++
				current := evt.Timestamp()
				duration := current.Sub(last)
				last = current
				durations = append(durations, duration)
				if len(durations) > b.count {
					durations = durations[1:]
				}
				total := 0 * time.Nanosecond
				low := 0x7FFFFFFFFFFFFFFF * time.Nanosecond
				high := 0 * time.Nanosecond
				for _, d := range durations {
					total += d
					if d < low {
						low = d
					}
					if d > high {
						high = d
					}
				}
				avg := total / time.Duration(len(durations))
				out.Emit(TopicRate, Rate{
					Time:             current,
					CountMatching:    countMatching,
					CountNonMatching: countNonMatching,
					Duration:         duration,
					High:             high,
					Low:              low,
					Average:          avg,
				})
			}
		}
	}
}

// EOF
