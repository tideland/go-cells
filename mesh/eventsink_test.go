// Tideland Go Cells - Mesh - Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
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

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"

	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestEventSinkPushPop verifies pushing and popping operations.
func TestEventSinkPushPop(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	max := 5
	sink := mesh.NewEventSink(max)
	assert.Length(sink, 0)

	evts := generateEvents(max)
	for i, evt := range evts {
		l := sink.Push(evt)
		assert.Equal(l, i+1)
	}
	assert.Length(sink, max)

	evts = generateEvents(1)
	evtA := evts[0]
	l := sink.Push(evtA)
	assert.Equal(l, max)
	assert.Length(sink, max)

	evtB, l := sink.Pop()
	assert.Equal(evtA, evtB)
	assert.Equal(l, max-1)

	for i := max - 1; i > 0; i-- {
		_, l = sink.Pop()
		assert.Equal(l, i-1)
	}
	assert.Length(sink, 0)
}

// TestEventSinkUnshiftShift verifies unshifting and shifting operations.
func TestEventSinkUnshiftShift(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	max := 5
	sink := mesh.NewEventSink(max)
	assert.Length(sink, 0)

	evts := generateEvents(max)
	for i, evt := range evts {
		l := sink.Unshift(evt)
		assert.Equal(l, i+1)
	}
	assert.Length(sink, max)

	evts = generateEvents(1)
	evtA := evts[0]
	l := sink.Unshift(evtA)
	assert.Equal(l, max)
	assert.Length(sink, max)

	evtB, l := sink.Shift()
	assert.Equal(evtA, evtB)
	assert.Equal(l, max-1)

	for i := max - 1; i > 0; i-- {
		_, l = sink.Shift()
		assert.Equal(l, i-1)
	}
	assert.Length(sink, 0)
}

// TestEventSinkFirstLastPeek verifies reading access to the sink.
func TestEventSinkFirstLastPeek(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	max := 5
	sink := mesh.NewEventSink(max)
	assert.Length(sink, 0)

	evts := generateEvents(max)
	evtFirst := evts[0]
	evtLast := evts[max-1]
	evtMid := evts[2]

	for _, evt := range evts {
		sink.Push(evt)
	}
	assert.Length(sink, max)

	first, ok := sink.First()
	assert.True(ok)
	last, ok := sink.Last()
	assert.True(ok)
	mid, ok := sink.Peek(2)
	assert.True(ok)

	assert.Equal(first, evtFirst)
	assert.Equal(last, evtLast)
	assert.Equal(mid, evtMid)

	assert.Length(sink, max)
}

// TestEventSinkDo verifies the iterating over a sink.
func TestEventSinkDo(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	// Do without error.
	evts := generateEvents(20)
	sink := mesh.NewEventSink(0, evts...)
	err := sink.Do(func(i int, evt *mesh.Event) error {
		assert.Equal(evt, evts[i])
		return nil
	})
	assert.NoError(err)

	// Do with error.
	err = sink.Do(func(i int, evt *mesh.Event) error {
		return errors.New("ouch")
	})
	assert.ErrorContains(err, "ouch")
}

// TestEventSinkFunctions verifies the functions on a sink reader.
func TestEventSinkFunctions(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	topics := []string{"a", "b", "b", "b", "a", "a", "b", "b", "a", "b"}
	sink := mesh.NewEventSink(0, generateTopicEvents(topics)...)

	// Filter events and reutrn error.
	evts, err := mesh.EventSinkFilter(sink, func(i int, evt *mesh.Event) (bool, error) {
		return evt.Topic() == "a", nil
	})
	assert.NoError(err)
	assert.Length(evts, 4)
	evts, err = mesh.EventSinkFilter(sink, func(i int, evt *mesh.Event) (bool, error) {
		return false, errors.New("ouch")
	})
	assert.ErrorContains(err, "ouch")
	assert.Length(evts, 0)

	// Match and mismatch events.
	ok, err := mesh.EventSinkMatch(sink, func(i int, evt *mesh.Event) (bool, error) {
		return evt.Topic() == "a" || evt.Topic() == "b", nil
	})
	assert.NoError(err)
	assert.OK(ok)
	ok, err = mesh.EventSinkMatch(sink, func(i int, evt *mesh.Event) (bool, error) {
		return evt.Topic() == "a" || evt.Topic() == "x", nil
	})
	assert.NoError(err)
	assert.False(ok)

	// Fold events.
	inject, err := mesh.NewEvent("counts", make(map[string]int))
	assert.NoError(err)
	facc, err := mesh.EventSinkFold(sink, inject, func(i int, acc, evt *mesh.Event) (*mesh.Event, error) {
		payload := make(map[string]int)
		err := acc.Payload(&payload)
		assert.NoError(err)
		payload[evt.Topic()]++
		return mesh.NewEvent(acc.Topic(), payload)
	})
	assert.NoError(err)
	payload := make(map[string]int)
	err = facc.Payload(&payload)
	assert.NoError(err)
	assert.Equal(payload["a"], 4)
	assert.Equal(payload["b"], 6)
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
