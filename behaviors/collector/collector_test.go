// Tideland Go Cells - Behaviors - Collector - Unit Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package collector_test // import "tideland.dev/go/cells/behaviors/collector"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"

	"tideland.dev/go/cells/behaviors/collector"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestSuccess verifies the successful usage of the collection behavior.
func TestSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	processor := func(r mesh.EventSinkReader) (*mesh.Event, error) {
		l := r.Len()
		return mesh.NewEvent("length", l)
	}
	behavior := collector.New(10, processor)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) {
			tbe.Push(evt)
		},
		func(tbe *mesh.TestbedEvaluator) {
			tbe.AssertRetry(func() bool { return tbe.Len() == 2 }, "assert 2 collected events: %v", tbe)
			evt, ok := tbe.Peek(0)
			tbe.Assert(ok, "first event missing")
			tbe.Assert(evt.Topic() == "length", "topic not equal 'length': %v", evt.Topic())
			var l int
			err := evt.Payload(&l)
			tbe.Assert(err == nil, "payload error not nil: %v", err)
			tbe.Assert(l == 10, "payload not 10: %v", l)
			evt, ok = tbe.Peek(1)
			tbe.Assert(ok, "second event missing")
			tbe.Assert(evt.Topic() == collector.TopicResetDone, "topic not equal 'reset-done': %v", evt.Topic())
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		for _, topic := range generator.Words(25) {
			out.Emit(topic)
		}
		out.Emit(collector.TopicProcess)
		out.Emit(collector.TopicReset)
	}, time.Second)
	assert.NoError(err)
}

// EOF
