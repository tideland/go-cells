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
	"strconv"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/cells/behaviors"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestCallbackBehavior tests the callback behavior.
func TestCallbackBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	count := 50
	countA := 0
	countB := 0
	callbackA := func(evt mesh.Event, out mesh.Emitter) error {
		return out.Emit("a")
	}
	callbackB := func(evt mesh.Event, out mesh.Emitter) error {
		return out.Emit("b")
	}
	behavior := behaviors.NewCallbackBehavior(callbackA, callbackB)
	eval := func(evt mesh.Event) (bool, error) {
		switch evt.Topic() {
		case "a":
			countA++
		case "b":
			countB++
		}
		return countA == countB && countA == count, nil
	}
	tb := mesh.NewTestbed(behavior, eval)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < count; i++ {
			topic := strconv.Itoa(i)
			out.Emit(topic)
		}
	}, time.Second)
	assert.NoError(err)
}

// EOF
