// Tideland Go Cells - Behaviors - Rate Window Evaluator
//
// Copyright (C) 2010-2022 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package ratewindow // import "tideland.dev/go/cells/behaviors/ratewindow"

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

// Behavior implements the rate window behavior. It can be used to check,
// if a number of wanted (matching) events happens in a defined timeframe.
// In case they are found all collected ones will be processed by a given
// function.
type Behavior struct {
	matches  RateWindowCriterion
	count    int
	duration time.Duration
	process  mesh.EventSinkProcessor
	sink     mesh.EventSink
}

// New creates a rate window behavior. Arguments are the function for
// checking the events, the number of expected matches, the duration
// and the processing function.
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
					continue
				}
				// Check matches and duration.
				b.sink.Push(evt)
				if b.sink.Len() == b.count {
					first, _ := b.sink.First()
					last, _ := b.sink.Last()
					difference := last.Timestamp().Sub(first.Timestamp())
					if difference <= b.duration {
						// We've got a burst, yeah!
						payload, err := b.process(b.sink)
						if err != nil {
							return err
						}
						out.Emit(TopicRateWindow, payload)
						b.sink.Clear()
					}
				}
			}
		}
	}
	return nil
}

// EOF
