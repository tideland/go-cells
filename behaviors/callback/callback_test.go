// Tideland Go Cells - Behaviors - Callback - Unit Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package callback_test // import "tideland.dev/go/cells/behaviors/callback"

//--------------------
// IMPORTS
//--------------------

import (
	"strconv"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/cells/behaviors/callback"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestSuccess verifies the successful call of callback functions.
func TestSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	count := 50
	countA := 0
	countB := 0
	callbackA := func(evt *mesh.Event, out mesh.Emitter) error {
		return out.Emit("a")
	}
	callbackB := func(evt *mesh.Event, out mesh.Emitter) error {
		return out.Emit("b")
	}
	behavior := callback.New(callbackA, callbackB)
	// Testing.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		switch evt.Topic() {
		case "a":
			countA++
		case "b":
			countB++
		}
		if countA == countB && countA == count {
			tbe.SetSuccess()
		}
		return nil
	}
	tb := mesh.NewTestbed(behavior, test)
	// Run tests.
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < count; i++ {
			topic := strconv.Itoa(i)
			out.Emit(topic)
		}
	}, time.Second)
	assert.NoError(err)
}

// EOF
