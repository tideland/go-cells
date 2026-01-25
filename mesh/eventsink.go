// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package mesh

import (
	"time"

	"tideland.dev/go/cells/mesh/internal"
)

// EventSinkDoFunc is used when looking over the collected events.
type EventSinkDoFunc = internal.EventSinkDoFunc

// EventSinkChanger can be used to write events into a sink or change
// it by reading and deleting.
type EventSinkChanger = internal.EventSinkChanger

// EventSinkReader can be used to read the events in a sink.
type EventSinkReader = internal.EventSinkReader

// EventSinkProcessor defines a function used to access the collected events
// in a sink for processing tasks. It only has reading access to the sink.
type EventSinkProcessor func(reader EventSinkReader) (any, error)

// EventSink combines changer and reader. It stores a number of ordered events by
// adding them at the end. To be used in behaviors for collecting sets of events
// and operate on them.
type EventSink = internal.EventSink

// NewEventSink creates a sink for events.
func NewEventSink(max int, evts ...*Event) EventSink {
	return internal.NewEventSink(max, evts...)
}

// EventSinkFilterFunc defines functions returning true for matching events.
type EventSinkFilterFunc func(i int, evt *Event) (bool, error)

// EventSinkFilter allows to retrieve a subset of events by running a filter function.
func EventSinkFilter(r EventSinkReader, filter EventSinkFilterFunc) ([]*Event, error) {
	var evts []*Event
	if derr := r.Do(func(i int, evt *Event) error {
		ok, err := filter(i, evt)
		if err != nil {
			return err
		}
		if ok {
			evts = append(evts, evt)
		}
		return nil
	}); derr != nil {
		return nil, derr
	}
	return evts, nil
}

// EventSinkMatch checks if all events match a passed filer.
func EventSinkMatch(r EventSinkReader, filter EventSinkFilterFunc) (bool, error) {
	matches := true
	if derr := r.Do(func(i int, evt *Event) error {
		ok, err := filter(i, evt)
		if err != nil {
			return err
		}
		matches = matches && ok
		return nil
	}); derr != nil {
		return false, derr
	}
	return matches, nil
}

// EventSinkFoldFunc defines functions for folding accelarator and event into a new event.
type EventSinkFoldFunc func(i int, acc, evt *Event) (*Event, error)

// EventSinkFold reduces (folds) the events of the sink to one.
func EventSinkFold(r EventSinkReader, inject *Event, fold EventSinkFoldFunc) (*Event, error) {
	var acc *Event = inject
	if derr := r.Do(func(i int, evt *Event) error {
		facc, err := fold(i, acc, evt)
		if err != nil {
			return err
		}
		acc = facc
		return nil
	}); derr != nil {
		return nil, derr
	}
	return acc, nil
}

// EventSinkDuration returns the duration between the first and the last event.
func EventSinkDuration(r EventSinkReader) time.Duration {
	first, fok := r.First()
	last, lok := r.Last()
	if fok == false || lok == false {
		return 0
	}
	return last.Timestamp().Sub(first.Timestamp())
}
