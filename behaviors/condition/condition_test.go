// Tideland Go Cells - Behaviors - Condition - Unit Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package condition_test // import "tideland.dev/go/cells/behaviors/condition"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"

	"tideland.dev/go/cells/behaviors/condition"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestSuccess verifies the successful scanning for conditions.
func TestSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "now"}
	tester := func(evt *mesh.Event) bool {
		return evt.Topic() == "now"
	}
	processor := func(cell mesh.Cell, evt *mesh.Event, out mesh.Emitter) error {
		topic := "found-" + evt.Topic()
		return out.Emit(topic)
	}
	behavior := condition.New(tester, processor)
	// Testing.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		if evt.Topic() == "found-now" {
			tbe.SignalSuccess()
		}
		return nil
	}
	// Run test.
	tb := mesh.NewTestbed(behavior, test)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 50; i++ {
			topic := generator.OneStringOf(topics...)
			out.Emit(topic)
		}
	}, time.Second)
	assert.NoError(err)
}

// EOF
