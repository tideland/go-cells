// Tideland Go Cells - Mesh - Tests
//
// Copyright (C) 2010-2026 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh_test // import "tideland.dev/go/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"testing"

	"tideland.dev/go/asserts/generators"
	"tideland.dev/go/asserts/verify"

	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestEventSinkPushPop verifies pushing and popping operations.
func TestEventSinkPushPop(t *testing.T) {

	max := 5
	sink := mesh.NewEventSink(max)
	verify.Length(t,sink, 0)

	evts := generateEvents(max)
	for i, evt := range evts {
		l := sink.Push(evt)
		verify.Equal(t,l, i+1)
	}
	verify.Length(t,sink, max)

	evts = generateEvents(1)
	evtA := evts[0]
	l := sink.Push(evtA)
	verify.Equal(t,l, max)
	verify.Length(t,sink, max)

	evtB, l := sink.Pop()
	verify.Equal(t,evtA, evtB)
	verify.Equal(t,l, max-1)

	for i := max - 1; i > 0; i-- {
		_, l = sink.Pop()
		verify.Equal(t,l, i-1)
	}
	verify.Length(t,sink, 0)
}

// TestEventSinkUnshiftShift verifies unshifting and shifting operations.
func TestEventSinkUnshiftShift(t *testing.T) {

	max := 5
	sink := mesh.NewEventSink(max)
	verify.Length(t,sink, 0)

	evts := generateEvents(max)
	for i, evt := range evts {
		l := sink.Unshift(evt)
		verify.Equal(t,l, i+1)
	}
	verify.Length(t,sink, max)

	evts = generateEvents(1)
	evtA := evts[0]
	l := sink.Unshift(evtA)
	verify.Equal(t,l, max)
	verify.Length(t,sink, max)

	evtB, l := sink.Shift()
	verify.Equal(t,evtA, evtB)
	verify.Equal(t,l, max-1)

	for i := max - 1; i > 0; i-- {
		_, l = sink.Shift()
		verify.Equal(t,l, i-1)
	}
	verify.Length(t,sink, 0)
}

// TestEventSinkFirstLastPeek verifies reading access to the sink.
func TestEventSinkFirstLastPeek(t *testing.T) {

	max := 5
	sink := mesh.NewEventSink(max)
	verify.Length(t,sink, 0)

	evts := generateEvents(max)
	evtFirst := evts[0]
	evtLast := evts[max-1]
	evtMid := evts[2]

	for _, evt := range evts {
		sink.Push(evt)
	}
	verify.Length(t,sink, max)

	first, ok := sink.First()
	verify.True(t,ok)
	last, ok := sink.Last()
	verify.True(t,ok)
	mid, ok := sink.Peek(2)
	verify.True(t,ok)

	verify.Equal(t,first, evtFirst)
	verify.Equal(t,last, evtLast)
	verify.Equal(t,mid, evtMid)

	verify.Length(t,sink, max)
}

// TestEventSinkDo verifies the iterating over a sink.
func TestEventSinkDo(t *testing.T) {

	// Do without error.
	evts := generateEvents(20)
	sink := mesh.NewEventSink(0, evts...)
	err := sink.Do(func(i int, evt *mesh.Event) error {
		verify.Equal(t,evt, evts[i])
		return nil
	})
	verify.NoError(t,err)

	// Do with error.
	err = sink.Do(func(i int, evt *mesh.Event) error {
		return errors.New("ouch")
	})
	verify.ErrorContains(t,err, "ouch")
}

// TestEventSinkFunctions verifies the functions on a sink reader.
func TestEventSinkFunctions(t *testing.T) {
	topics := []string{"a", "b", "b", "b", "a", "a", "b", "b", "a", "b"}
	sink := mesh.NewEventSink(0, generateTopicEvents(topics)...)

	// Filter events and reutrn error.
	evts, err := mesh.EventSinkFilter(sink, func(i int, evt *mesh.Event) (bool, error) {
		return evt.Topic() == "a", nil
	})
	verify.NoError(t,err)
	verify.Length(t,evts, 4)
	evts, err = mesh.EventSinkFilter(sink, func(i int, evt *mesh.Event) (bool, error) {
		return false, errors.New("ouch")
	})
	verify.ErrorContains(t,err, "ouch")
	verify.Length(t,evts, 0)

	// Match and mismatch events.
	ok, err := mesh.EventSinkMatch(sink, func(i int, evt *mesh.Event) (bool, error) {
		return evt.Topic() == "a" || evt.Topic() == "b", nil
	})
	verify.NoError(t,err)
	verify.True(t,ok)
	ok, err = mesh.EventSinkMatch(sink, func(i int, evt *mesh.Event) (bool, error) {
		return evt.Topic() == "a" || evt.Topic() == "x", nil
	})
	verify.NoError(t,err)
	verify.False(t,ok)

	// Fold events.
	inject, err := mesh.NewEvent("counts", make(map[string]int))
	verify.NoError(t,err)
	facc, err := mesh.EventSinkFold(sink, inject, func(i int, acc, evt *mesh.Event) (*mesh.Event, error) {
		payload := make(map[string]int)
		err := acc.Payload(&payload)
		verify.NoError(t,err)
		payload[evt.Topic()]++
		return mesh.NewEvent(acc.Topic(), payload)
	})
	verify.NoError(t,err)
	payload := make(map[string]int)
	err = facc.Payload(&payload)
	verify.NoError(t,err)
	verify.Equal(t,payload["a"], 4)
	verify.Equal(t,payload["b"], 6)
}

//--------------------
// HELPER
//--------------------

// generateEvents generates a number of events for tests.
func generateEvents(count int) []*mesh.Event {
	generator := generators.New(generators.FixedRand())
	topics := generator.Words(count)
	return generateTopicEvents(topics)
}

// generateTopicEvents generates a number of events for tests
// based on topics.
func generateTopicEvents(topics []string) []*mesh.Event {
	evts := []*mesh.Event{}
	for _, topic := range topics {
		evt, _ := mesh.NewEvent(topic)
		evts = append(evts, evt)
	}
	return evts
}

// EOF
