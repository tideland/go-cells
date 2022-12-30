// Tideland Go Cells - Behaviors - Rate Window Evaluator - Unit Tests
//
// Copyright (C) 2010-2022 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package ratewindow_test // import "tideland.dev/go/cells/behaviors/ratewindow"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"

	"tideland.dev/go/cells/behaviors/ratewindow"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestSuccess verifies the successful finding and processing of matching events.
func TestSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	matcher := func(evt *mesh.Event) (bool, error) {
		var payload int
		err := evt.Payload(&payload)
		assert.NoError(err)
		return payload > 5, nil
	}
	processor := func(reader mesh.EventSinkReader) (interface{}, error) {
		var count int
		var sum int
		doer := func(i int, evt *mesh.Event) error {
			var payload int
			err := evt.Payload(&payload)
			assert.NoError(err)
			count += 1
			sum += payload
			return nil
		}
		err := reader.Do(doer)
		return sum / count, err
	}
	behavior := ratewindow.New(matcher, 50, time.Second, processor)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			tbe.WaitFor(func() bool { return tbe.Len() > 0 })
			tbe.Do(func(i int, evt *mesh.Event) error {
				var payload int
				err := evt.Payload(&payload)
				tbe.Assert(err == nil, "invalid payload: payload is no int")
				tbe.Assert(payload > 5, "invalid payload: payload is <= 5")
				return nil
			})
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 10000; i++ {
			topic := "int"
			payload := generator.Int(0, 10)
			out.Emit(topic, payload)
		}
	}, 10*time.Second)
	assert.NoError(err)
}

// EOF
