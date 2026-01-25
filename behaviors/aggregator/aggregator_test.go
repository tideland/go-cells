// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package aggregator_test

import (
	"strconv"
	"testing"
	"time"

	"tideland.dev/go/asserts/verify"

	"tideland.dev/go/cells/behaviors/aggregator"
	"tideland.dev/go/cells/mesh"
)

// TestAggregatorBehavior tests the aggregator behavior.
func TestAggregatorBehavior(t *testing.T) {
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
	verify.NoError(t, err)
}
