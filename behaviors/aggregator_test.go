// Tideland Go Cells - Behaviors
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors_test // import "tideland.dev/go/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"strconv"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/cells/behaviors"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestAggregatorBehavior tests the aggregator behavior.
func TestAggregatorBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	count := 50
	aggregator := func(aggregate interface{}, evt *mesh.Event) (interface{}, error) {
		words := aggregate.(map[string]bool)
		words[evt.Topic()] = true
		return words, nil
	}
	behavior := behaviors.NewAggregatorBehavior(map[string]bool{}, aggregator)
	// Test evaluation.
	eval := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		tbe.Push(evt)
		if tbe.Len() == 2 {
			// Check first for aggregated.
			evtA, _ := tbe.First()
			if evtA.Topic() != behaviors.TopicAggregated {
				tbe.SetFail("topic not 'aggregated': %s", evtA.Topic())
				return nil
			}
			var words map[string]bool
			if err := evtA.Payload(&words); err != nil {
				return err
			}
			if len(words) != count {
				tbe.SetFail("invalid length of words: %d", len(words))
				return nil
			}
			// Check second for resetted.
			evtB, _ := tbe.First()
			if evtB.Topic() != behaviors.TopicResetted {
				tbe.SetFail("topic not 'resetted': %s", evtB.Topic())
				return nil
			}
			tbe.SetSuccess()
		}
		return nil
	}
	// Run tests.
	tb := mesh.NewTestbed(behavior, eval)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < count; i++ {
			topic := strconv.Itoa(i)
			out.Emit(topic)
		}
		out.Emit(behaviors.TopicAggregate)
		out.Emit(behaviors.TopicReset)
	}, time.Second)
	assert.NoError(err)
}

// EOF
