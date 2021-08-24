// Tideland Go Cells - Behaviors - Aggregator - Unit Test
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package aggregator_test // import "tideland.dev/go/cells/behaviors/aggregator"

//--------------------
// IMPORTS
//--------------------

import (
	"strconv"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/cells/behaviors/aggregator"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestAggregatorBehavior tests the aggregator behavior.
func TestAggregatorBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	count := 50
	initializer := func() interface{} {
		return map[string]bool{
			"initialized": true,
		}
	}
	aggregatorFunc := func(aggregated interface{}, evt *mesh.Event) (interface{}, error) {
		words := aggregated.(map[string]bool)
		words[evt.Topic()] = true
		return words, nil
	}
	behavior := aggregator.New(initializer, aggregatorFunc)
	// Test evaluation.
	eval := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		tbe.Push(evt)
		switch evt.Topic() {
		case aggregator.TopicAggregateDone:
			var words map[string]bool
			if err := evt.Payload(&words); err != nil {
				return err
			}
			if len(words) != count+1 {
				tbe.SetFail("invalid length of aggregated words: %d", len(words))
				return nil
			}
		case aggregator.TopicResetDone:
			var words map[string]bool
			if err := evt.Payload(&words); err != nil {
				return err
			}
			if len(words) != 1 {
				tbe.SetFail("invalid length of resetted words: %d", len(words))
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
		out.Emit(aggregator.TopicAggregate)
		out.Emit(aggregator.TopicReset)
	}, time.Second)
	assert.NoError(err)
}

// EOF
