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
	initializer := func() any {
		return map[string]bool{
			"initialized": true,
		}
	}
	aggregatorFunc := func(aggregated any, evt *mesh.Event) (any, error) {
		words := aggregated.(map[string]bool)
		words[evt.Topic()] = true
		return words, nil
	}
	behavior := aggregator.New(initializer, aggregatorFunc)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			tbe.WaitFor(func() bool { return tbe.Len() == 2 })
			evt, ok := tbe.First()
			tbe.Assert(ok, "cannot retrieve first event")
			tbe.Assert(evt.Topic() == aggregator.TopicAggregateDone, "first event is no aggregate done event: %v", evt)
			var fstWords map[string]bool
			tbe.Assert(evt.Payload(&fstWords) == nil, "retrieving of first payload failed")
			tbe.Assert(len(fstWords) == count+1, "invalid length of aggregated words: %d", len(fstWords))
			evt, ok = tbe.Last()
			tbe.Assert(ok, "cannot retrieve last event")
			tbe.Assert(evt.Topic() == aggregator.TopicResetDone, "last event is no reset done event: %v", evt)
			var sndWords map[string]bool
			tbe.Assert(evt.Payload(&sndWords) == nil, "retrieving of first payload failed")
			tbe.Assert(len(sndWords) == 1, "invalid length of resetted words: %d", len(sndWords))
		},
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
