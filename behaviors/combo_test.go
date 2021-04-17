// Tideland Go Cells - Behaviors - Unit Tests
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
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"

	"tideland.dev/go/cells/behaviors"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestComboBehavior tests the combo behavior.
func TestComboBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	topics := generator.Words(50)
	wanted := "test-topic"
	matcher := func(r mesh.EventSinkReader) (behaviors.CriterionMatch, interface{}, error) {
		// Matcher tries to find the wanted topic twice. When found twice
		// the distance will be returned.
		var found []int
		if err := r.Do(func(i int, evt *mesh.Event) error {
			if evt.Topic() == wanted {
				found = append(found, i)
			}
			return nil
		}); err != nil {
			return behaviors.CriterionError, nil, err
		}
		// Check if found and where.
		switch {
		case len(found) == 2:
			return behaviors.CriterionDone, found[1] - found[0], nil
		case len(found) == 1:
			if found[0] == 0 {
				return behaviors.CriterionKeep, nil, nil
			}
		}
		return behaviors.CriterionDropFirst, nil, nil
	}
	behavior := behaviors.NewComboBehavior(matcher)
	eval := func(evt *mesh.Event) (bool, error) {
		switch evt.Topic() {
		case behaviors.TopicCriterionDone:
			var distance int
			err := evt.Payload(&distance)
			return distance == 50, err
		case mesh.TopicTestbedTerminated:
			return true, nil
		}
		return false, nil
	}
	// Run test.
	tb := mesh.NewTestbed(behavior, eval)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 100; i++ {
			var topic string
			if i == 25 || i == 75 {
				topic = wanted
			} else {
				topic = generator.OneStringOf(topics...)
			}
			err := out.Emit(topic)
			assert.NoError(err)
		}
	}, time.Second)
	assert.NoError(err)
}

// EOF
