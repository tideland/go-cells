// Tideland Go Cells - Behaviors - Countdown - Unit Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package countdown_test // import "tideland.dev/go/cells/behaviors/countdown"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"

	"tideland.dev/go/cells/behaviors/countdown"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestSuccess verifies the countdown of events.
func TestSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	count := 10
	zeroer := func(r mesh.EventSinkReader) (*mesh.Event, error) {
		if r.Len() == count {
			return mesh.NewEvent("length-ok")
		}
		return nil, nil
	}
	behavior := countdown.New(count, zeroer)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			tbe.AssertRetry(func() bool { return tbe.Len() == 1 }, "not yet all events emitted")
			evt, ok := tbe.First()
			tbe.Assert(ok, "no first event")
			tbe.Assert(evt.Topic() == "length-ok", "topic is not 'length-ok': %v", evt.Topic())
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < count; i++ {
			topic := generator.OneStringOf(topics...)
			out.Emit(topic)
		}
	}, time.Second)
	assert.NoError(err)
}

// EOF
