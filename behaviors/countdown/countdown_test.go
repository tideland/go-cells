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
	// Test evaluation.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		if evt.Topic() == "length-ok" {
			tbe.SetSuccess()
		}
		return nil
	}
	// Run test.
	tb := mesh.NewTestbed(behavior, test)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < count; i++ {
			topic := generator.OneStringOf(topics...)
			out.Emit(topic)
		}
	}, time.Second)
	assert.NoError(err)
}

// EOF
