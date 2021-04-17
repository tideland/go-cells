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
	eval := func(evt *mesh.Event) (bool, error) {
		switch evt.Topic() {
		case behaviors.TopicResetted:
			return true, nil
		case behaviors.TopicAggregated:
			var words map[string]bool
			err := evt.Payload(&words)
			assert.NoError(err)
			assert.Length(words, count)
		}
		return false, nil
	}
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
