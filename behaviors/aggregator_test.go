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
		var words map[string]bool
		if !evt.HasPayload() {
			words = map[string]bool{}
		} else {
			words = aggregate.(map[string]bool)
		}
		words[evt.Topic()] = true
		return words, nil
	}
	behavior := behaviors.NewAggregatorBehavior(aggregator)
	// Test evaluation.
	eval := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		tbe.Push(evt)
		switch evt.Topic() {
		case behaviors.TopicAggregated:
			var words map[string]bool
			if err := evt.Payload(&words); err != nil {
				return err
			}
			if len(words) != count+1 {
				tbe.SetFail("invalid length of aggregated words: %d", len(words))
				return nil
			}
		case behaviors.TopicResetted:
			if evt.HasPayload() {
				tbe.SetFail("event has illegal payload")
				return nil
			}
		}
		if tbe.Len() == 2 {
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
