// Tideland Go Cells - Behaviors - Rate Window Evaluator
//
// Copyright (C) 2010-2022 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package rateevaluator // import "tideland.dev/go/cells/behaviors/ratewindow"

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
	// TopicReset tells the cell to reset its collected events.
	TopicReset = "reset!"

	// TopicRateWindow signals a detected event rate window.
	TopicRateWindow = "rate-window"
)

//--------------------
// HELPER
//--------------------

// RateWindowCriterion is used by the rate window behavior and has to return
// true, if the passed event matches a criterion for rate window measuring.
type RateWindowCriterion func(evt *mesh.Event) (bool, error)

//--------------------
// RATE WINDOW BEHAVIOR
//--------------------

// Behavior implements the rate window behavior.
type Behavior struct {
	matches  RateWindowCriterion
	count    int
	duration time.Duration
	process  mesh.EventSinkProcessor
	sink     mesh.EventSink
}

// New creates an event rate window behavior. It checks if an event
// matches the passed criterion. If count events match during
// duration the process function is called. Its returned payload is
// emitted as new event with topic "rate-window". A received "reset" as
// topic resets the collected matches.
func New(
	matches RateWindowCriterion,
	count int,
	duration time.Duration,
	process mesh.EventSinkProcessor) *Behavior {
	return &Behavior{
		matches:  matches,
		count:    count,
		duration: duration,
		process:  process,
		sink:     mesh.NewEventSink(count),
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
			default:
				// CHeck if the event matches.
				ok, err := b.matches(evt)
				if err != nil {
					return err
				}
				if !ok {
					return nil
				}
				// Check matches and duration.
				b.sink.Push(evt)
				if b.sink.Len() == b.count {
					first, _ := b.sink.First()
					last, _ := b.sink.Last()
					difference := last.Timestamp().Sub(first.Timestamp())
					if difference <= b.duration {
						// We've got a burst!
						payload, err := b.process(b.sink)
						if err != nil {
							return err
						}
						out.Emit(TopicRateWindow, payload)
					}
					b.sink.Shift()
				}
			}
		}
	}
	return nil
}

// EOF
