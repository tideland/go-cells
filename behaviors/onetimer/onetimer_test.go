// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package onetimer_test

import (
	"testing"
	"time"

	"tideland.dev/go/asserts/verify"

	"tideland.dev/go/cells/behaviors/onetimer"
	"tideland.dev/go/cells/mesh"
)

// TestSuccess verifies the
func TestSuccess(t *testing.T) {
	count := 0
	oneTime := func(evt *mesh.Event, out mesh.Emitter) error {
		count++
		verify.True(t, count < 2)
		out.EmitEvent(evt)
		return nil
	}
	behavior := onetimer.New(oneTime)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			// Test three times as those could be emitted later.
			tbe.AssertRetry(func() bool { return tbe.Len() == 1 }, "invalid number of emitted events: %v", tbe)
			tbe.AssertRetry(func() bool { return tbe.Len() == 1 }, "invalid number of emitted events: %v", tbe)
			tbe.AssertRetry(func() bool { return tbe.Len() == 1 }, "invalid number of emitted events: %v", tbe)
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("a")
		out.Emit("b")
		out.Emit("c")
	}, time.Second)
	verify.NoError(t, err)
}
