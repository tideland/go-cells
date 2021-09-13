// Tideland Go Cells - Behaviors - Unit Tests - Once
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package onetimer_test // import "tideland.dev/go/cells/behaviors/onetimer"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/cells/behaviors/onetimer"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestSuccess verifies the
func TestSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	count := 0
	oneTime := func(evt *mesh.Event, out mesh.Emitter) error {
		count++
		assert.True(count < 2)
		out.EmitEvent(evt)
		return nil
	}
	behavior := onetimer.New(oneTime)
	// Testing.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		switch evt.Topic() {
		case "a":
			tbe.SetSuccess()
		case "b", "c":
			tbe.SetFail("received invalid events")
		}
		return nil
	}
	// Run test.
	tb := mesh.NewTestbed(behavior, test)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("a")
		out.Emit("b")
		out.Emit("c")
	}, time.Second)
	assert.NoError(err)
}

// EOF
