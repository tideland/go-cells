// Tideland Go Cells - Behaviors - Broadcaster - Unit Test
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package broadcaster_test // import "tideland.dev/go/cells/behaviors/broadcaster"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/cells/behaviors/broadcaster"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestSuccess verifies the successfull broadcasting.
func TestSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	behavior := broadcaster.New()
	// Test behavior, pretty simple as the evaluator tests all emitted events.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		switch evt.Topic() {
		case "done":
			tbe.SignalSuccess()
		}
		return nil
	}
	// Run tests.
	tb := mesh.NewTestbed(behavior, test)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("one")
		out.Emit("two")
		out.Emit("three")
		out.Emit("done")
	}, time.Second)
	assert.NoError(err)
}

// EOF
