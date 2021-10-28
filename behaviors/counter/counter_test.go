// Tideland Go Cells - Behaviors - Counter - Unit Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package counter_test // import "tideland.dev/go/cells/behaviors/counter"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"

	"tideland.dev/go/cells/behaviors/counter"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestSuccess tests the successful counting of events and resetting the counters.
func TestSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	topics := []string{"alpha", "bravo", "charly", "delta", "echo"}
	counteval := func(evt *mesh.Event) ([]string, error) {
		var counters []string
		for _, r := range evt.Topic() {
			counter := fmt.Sprintf("%v", r)
			counters = append(counters, counter)
		}
		return counters, nil
	}
	behavior := counter.New(counteval)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			clen := func(evt *mesh.Event) int {
				var counters map[string]int
				if err := evt.Payload(&counters); err != nil {
					return -1
				}
				return len(counters)
			}
			tbe.AssertRetry(func() bool { return tbe.Len() == 3 }, "wait for three emitted events failed")
			evt, ok := tbe.Peek(0)
			tbe.Assert(ok, "retrieving first event failed")
			tbe.Assert(evt.Topic() == counter.TopicCountersDone, "first topic failed: %v", evt)
			tbe.Assert(clen(evt) == 13, "not 13 counters after first run: %v", evt)
			evt, ok = tbe.Peek(1)
			tbe.Assert(ok, "retrieving second event failed")
			tbe.Assert(evt.Topic() == counter.TopicResetDone, "second topic failed: %v", evt)
			evt, ok = tbe.Peek(2)
			tbe.Assert(ok, "retrieving third event failed")
			tbe.Assert(evt.Topic() == counter.TopicCountersDone, "third topic failed: %v", evt)
			tbe.Assert(clen(evt) == 0, "not 0 counters after third run: %v", evt)
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 50; i++ {
			topic := generator.OneStringOf(topics...)
			out.Emit(topic)
		}
		// Retrieve, erase, and retrieve again.
		out.Emit(counter.TopicCounters)
		out.Emit(counter.TopicReset)
		out.Emit(counter.TopicCounters)
	}, time.Second)
	assert.NoError(err)
}

// EOF
