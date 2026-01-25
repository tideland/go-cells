// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package callback_test

import (
	"strconv"
	"testing"
	"time"

	"tideland.dev/go/asserts/verify"

	"tideland.dev/go/cells/behaviors/callback"
	"tideland.dev/go/cells/mesh"
)

// TestSuccess verifies the successful call of callback functions.
func TestSuccess(t *testing.T) {
	count := 50
	callbackA := func(evt *mesh.Event, out mesh.Emitter) error {
		return out.Emit("a")
	}
	callbackB := func(evt *mesh.Event, out mesh.Emitter) error {
		return out.Emit("b")
	}
	behavior := callback.New(callbackA, callbackB)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			tbe.WaitFor(func() bool { return tbe.Len() == 100 })
			countA := 0
			countB := 0
			tbe.Do(func(i int, evt *mesh.Event) error {
				switch evt.Topic() {
				case "a":
					countA++
				case "b":
					countB++
				}
				return nil
			})
			tbe.Assert(countA == count, "counter A is wrong: %d", countA)
			tbe.Assert(countB == count, "counter B is wrong: %d", countB)
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < count; i++ {
			topic := strconv.Itoa(i)
			out.Emit(topic)
		}
	}, time.Second)
	verify.NoError(t, err)
}
