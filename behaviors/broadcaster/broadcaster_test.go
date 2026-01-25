// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package broadcaster_test

import (
	"testing"
	"time"

	"tideland.dev/go/asserts/verify"

	"tideland.dev/go/cells/behaviors/broadcaster"
	"tideland.dev/go/cells/mesh"
)

// TestSuccess verifies the successfull broadcasting.
func TestSuccess(t *testing.T) {
	behavior := broadcaster.New()
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			tbe.AssertRetry(func() bool { return tbe.Len() == 3 }, "broadcasted events not 3: %v", tbe)
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("one")
		out.Emit("two")
		out.Emit("three")
	}, time.Second)
	verify.NoError(t, err)
}
