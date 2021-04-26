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
	"fmt"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/cells/behaviors"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestBroadcasterBehavior tests the broadcaster behavior.
func TestBroadcasterBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	behavior := behaviors.NewBroadcasterBehavior()
	values := map[string]bool{
		"one":                true,
		"two":                true,
		"three":              true,
		"done":               true,
		"testbed-terminated": true,
	}
	// Test evaluation.
	eval := func(tbctx *mesh.TestbedContext, evt *mesh.Event) error {
		knowsValue := values[evt.Topic()]
		if !knowsValue {
			return fmt.Errorf("unknown topic: %s", evt.Topic())
		}
		if evt.Topic() == "done" {
			tbctx.SetSuccess()
		}
		return nil
	}
	// Run tests.
	tb := mesh.NewTestbed(behavior, eval)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("one")
		out.Emit("two")
		out.Emit("three")
		out.Emit("done")
	}, time.Second)
	assert.NoError(err)
}

// EOF
