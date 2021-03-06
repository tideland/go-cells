// Tideland Go Cells - Behaviors - Combo - Unit Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package combo_test // import "tideland.dev/go/cells/behaviors/combo"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"

	"tideland.dev/go/cells/behaviors/combo"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestSuccess verifies the successful double finding of a topic.
func TestSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	topics := generator.Words(50)
	wanted := "test-topic"
	matcher := func(r mesh.EventSinkReader) (combo.CriterionMatch, interface{}, error) {
		// Matcher tries to find the wanted topic twice. When found twice
		// the distance will be returned.
		var found []int
		if err := r.Do(func(i int, evt *mesh.Event) error {
			if evt.Topic() == wanted {
				found = append(found, i)
			}
			return nil
		}); err != nil {
			return combo.CriterionError, nil, err
		}
		// Check if found and where.
		switch {
		case len(found) == 2:
			return combo.CriterionDone, found[1] - found[0], nil
		case len(found) == 1:
			if found[0] == 0 {
				return combo.CriterionKeep, nil, nil
			}
		}
		return combo.CriterionDropFirst, nil, nil
	}
	behavior := combo.New(matcher)
	// Run test.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			tbe.AssertRetry(func() bool { return tbe.Len() > 0 }, "got at least one combo")
			err := tbe.Do(func(i int, evt *mesh.Event) error {
				tbe.Assert(evt.Topic() == combo.TopicCriterionDone, "topic %d does not show done combo", i)
				var distance int
				if err := evt.Payload(&distance); err != nil {
					return err
				}
				tbe.Assert(distance == 50, "invalid distance of combo: %d", distance)

				return nil
			})
			tbe.Assert(err == nil, "testing of collected events had error: %v", err)
		},
	)
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
