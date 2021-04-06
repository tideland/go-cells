// Tideland Go Cells - Behaviors - Unit Tests
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
	"tideland.dev/go/audit/generators"

	"tideland.dev/go/cells/behaviors"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestCollectorBehavior tests the collector behavior.
func TestCollectorBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	processor := func(r mesh.EventSinkReader) (mesh.Event, error) {
		l := r.Len()
		return mesh.NewEvent("length", l)
	}
	behavior := behaviors.NewCollectorBehavior(10, processor)
	eval := func(evt mesh.Event) (bool, error) {
		switch evt.Topic() {
		case "length":
			var l int
			err := evt.Payload(&l)
			assert.NoError(err)
			assert.Equal(l, 10)
			return false, nil
		case behaviors.TopicResetted:
			return true, nil
		}
		return false, nil
	}
	tb := mesh.NewTestbed(behavior, eval)
	err := tb.Go(func(out mesh.Emitter) {
		for _, topic := range generator.Words(25) {
			out.Emit(topic)
		}
		out.Emit(behaviors.TopicProcess)
		out.Emit(behaviors.TopicReset)
	}, time.Second)
	assert.NoError(err)
}

// EOF
