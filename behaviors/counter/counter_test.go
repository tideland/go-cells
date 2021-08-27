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
	// Testing.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		switch evt.Topic() {
		case counter.TopicCountersDone:
			var counters map[string]int
			err := evt.Payload(&counters)
			if err != nil {
				return err
			}
			l := len(counters)
			// 13 characters = 13 counters and 0 after a reset.
			if l == 13 || l == 0 {
				tbe.SetSuccess()
			}
			return nil
		default:
			return nil
		}
	}
	// Run test.
	tb := mesh.NewTestbed(behavior, test)
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
