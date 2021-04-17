// Tideland Go Cells - Mesh
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh // import "tideland.dev/go/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"time"
)

//--------------------
// EVENT SINK READER
//--------------------

// EventSinkDoFunc is used when looking over the collected events.
type EventSinkDoFunc func(i int, evt *Event) error

// EventSinkReader can be used to read the events in a sink. It is a
// specialized subfunctionality of the event sink.
type EventSinkReader interface {
	// Len returns the number of stored events.
	Len() int

	// PeekFirst returns the first of the collected events.
	First() (*Event, bool)

	// PeekLast returns the last of the collected event datas.
	Last() (*Event, bool)

	// PeekAt returns an event at a given index and true if it
	// exists, otherwise nil and false.
	Peek(index int) (*Event, bool)

	// Do iterates over all collected events.
	Do(do EventSinkDoFunc) error
}

//--------------------
// EVENT SINK
//--------------------

// EventSink stores a number of ordered events by adding them at the end. To
// be used in behaviors for collecting sets of events and operate on them.
type EventSink struct {
	max    int
	events []*Event
}

// NewEventSink creates a sink for events.
func NewEventSink(max int, evts ...*Event) *EventSink {
	s := &EventSink{
		max: max,
	}
	if max > 0 && len(evts) > max {
		s.events = append(s.events, evts[len(evts)-max:]...)
	} else {
		s.events = append(s.events, evts...)
	}
	return s
}

// Push adds an event to the end of the sink.
func (s *EventSink) Push(evt *Event) int {
	s.events = append(s.events, evt)
	if s.max > 0 && len(s.events) > s.max {
		s.events = s.events[1:]
	}
	return len(s.events)
}

// Pop retrieves and removes the last event from the sink
// and also returns the new length.
func (s *EventSink) Pop() (*Event, int) {
	if len(s.events) == 0 {
		return nil, 0
	}
	l := len(s.events) - 1
	evt := s.events[l]
	s.events = s.events[:l]
	return evt, l
}

// Unshift adds an event to the begin of the sink.
func (s *EventSink) Unshift(evt *Event) int {
	s.events = append([]*Event{evt}, s.events...)
	if s.max > 0 && len(s.events) > s.max {
		s.events = s.events[:len(s.events)-1]
	}
	return len(s.events)
}

// Shift returns and removes the first event of the sink
// and also returns the new length.
func (s *EventSink) Shift() (*Event, int) {
	if len(s.events) == 0 {
		return nil, 0
	}
	l := len(s.events) - 1
	evt := s.events[0]
	s.events = s.events[1:]
	return evt, l
}

// First allows a look at the first event of the sink if it
// exists. Otherwise nil and false will be returned.
func (s *EventSink) First() (*Event, bool) {
	if len(s.events) < 1 {
		return nil, false
	}
	return s.events[0], true
}

// Last allows a look at the last event of the sink if it
// exists. Otherwise nil and false will be returned.
func (s *EventSink) Last() (*Event, bool) {
	if len(s.events) < 1 {
		return nil, false
	}
	return s.events[len(s.events)-1], true
}

// Peek allows a look at the indexed event of the sink if it
// exists. Otherwise nil and false will be returned.
func (s *EventSink) Peek(index int) (*Event, bool) {
	if index < 0 || index > len(s.events)-1 {
		return nil, false
	}
	return s.events[index], true
}

// Clear removes all collected events.
func (s *EventSink) Clear() {
	s.events = nil
}

// Len returns the number of events in the sink.
func (s *EventSink) Len() int {
	return len(s.events)
}

// Do allows to iterate over all events of the sink and perform a
// function.
func (s *EventSink) Do(do EventSinkDoFunc) error {
	for i, evt := range s.events {
		if err := do(i, evt); err != nil {
			return err
		}
	}
	return nil
}

//--------------------
// EVENT SINK FUNCTIONS
//--------------------

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

// EOF
