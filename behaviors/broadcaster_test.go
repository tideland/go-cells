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
	topics := make(map[string]bool)
	eval := func(evt mesh.Event) (bool, error) {
		if evt.Topic() == "done" {
			return true, nil
		}
		topics[evt.Topic()] = true
		return false, nil
	}
	tb := mesh.NewTestbed(behavior, eval)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("one")
		out.Emit("two")
		out.Emit("three")
		out.Emit("done")
	}, time.Second)
	assert.NoError(err)

	assert.True(topics["one"])
	assert.True(topics["two"])
	assert.True(topics["three"])
	assert.False(topics["done"])
}

// EOF
