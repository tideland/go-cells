// Tideland Go Cells - Mesh - Internal
//
// Copyright (C) 2010-2026 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package internal // import "tideland.dev/go/cells/mesh/internal"

//--------------------
// IMPORTS
//--------------------

import (
	"strings"
	"sync"
)

//--------------------
// EVENT SINK
//--------------------

// EventSink stores a number of ordered events by adding them at the end. To
// be used in behaviors for collecting sets of events and operate on them.
type eventSinkImpl struct {
	mux    sync.RWMutex
	max    int
	events []*Event
}

// NewEventSink creates a sink for events.
func NewEventSink(max int, evts ...*Event) *eventSinkImpl {
	s := &eventSinkImpl{
		max: max,
	}
	if max > 0 && len(evts) > max {
		s.events = append(s.events, evts[len(evts)-max:]...)
	} else {
		s.events = append(s.events, evts...)
	}
	return s
}

// Push adds an event to the end of the sink and returns the new size.
func (s *eventSinkImpl) Push(evt *Event) int {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.events = append(s.events, evt)
	if s.max > 0 && len(s.events) > s.max {
		s.events = s.events[1:]
	}
	return len(s.events)
}

// Pop retrieves and removes the last event from the sink
// and also returns the new length.
func (s *eventSinkImpl) Pop() (*Event, int) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if len(s.events) == 0 {
		return nil, 0
	}
	l := len(s.events) - 1
	evt := s.events[l]
	s.events = s.events[:l]
	return evt, l
}

// Unshift adds an event to the begin of the sink and returns the new size.
func (s *eventSinkImpl) Unshift(evt *Event) int {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.events = append([]*Event{evt}, s.events...)
	if s.max > 0 && len(s.events) > s.max {
		s.events = s.events[:len(s.events)-1]
	}
	return len(s.events)
}

// Shift returns and removes the first event of the sink
// and also returns the new length.
func (s *eventSinkImpl) Shift() (*Event, int) {
	s.mux.Lock()
	defer s.mux.Unlock()
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
func (s *eventSinkImpl) First() (*Event, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	if len(s.events) < 1 {
		return nil, false
	}
	return s.events[0], true
}

// Last allows a look at the last event of the sink if it
// exists. Otherwise nil and false will be returned.
func (s *eventSinkImpl) Last() (*Event, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	if len(s.events) < 1 {
		return nil, false
	}
	return s.events[len(s.events)-1], true
}

// Peek allows a look at the indexed event of the sink if it
// exists. Otherwise nil and false will be returned.
func (s *eventSinkImpl) Peek(index int) (*Event, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	if index < 0 || index > len(s.events)-1 {
		return nil, false
	}
	return s.events[index], true
}

// Clear removes all collected events.
func (s *eventSinkImpl) Clear() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.events = nil
}

// Len returns the number of events in the sink.
func (s *eventSinkImpl) Len() int {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return len(s.events)
}

// Do allows to iterate over all events of the sink and perform a
// function.
func (s *eventSinkImpl) Do(do EventSinkDoFunc) error {
	s.mux.RLock()
	defer s.mux.RUnlock()
	for i, evt := range s.events {
		if err := do(i, evt); err != nil {
			return err
		}
	}
	return nil
}

// String prints the topics inside the sink.
func (s *eventSinkImpl) String() string {
	s.mux.RLock()
	defer s.mux.RUnlock()
	var topics []string
	s.Do(func(i int, evt *Event) error {
		topics = append(topics, "\""+evt.Topic()+"\"")
		return nil
	})
	return "[" + strings.Join(topics, " ") + "]"
}

// EOF
