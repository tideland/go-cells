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
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) {
			tbe.Push(evt)
			switch evt.Topic() {
			case aggregator.TopicAggregateDone:
				var words map[string]bool
				if err := evt.Payload(&words); err != nil {
					tbe.SignalError(err)
				}
				tbe.Assert(len(words) == count+1, "invalid length of aggregated words: %d", len(words))
			case aggregator.TopicResetDone:
				var words map[string]bool
				if err := evt.Payload(&words); err != nil {
					tbe.SignalError(err)
				}
				tbe.Assert(len(words) == 1, "invalid length of resetted words: %d", len(words))
			}
		},
		nil,
	)
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
