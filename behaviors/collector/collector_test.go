// Tideland Go Cells - Behaviors - Collector - Unit Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package collector_test // import "tideland.dev/go/cells/behaviors/collector"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"

	"tideland.dev/go/cells/behaviors"
	"tideland.dev/go/cells/behaviors/collector"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestSuccess verifies the successful isage of the collection behavior.
func TestSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	processor := func(r mesh.EventSinkReader) (*mesh.Event, error) {
		l := r.Len()
		return mesh.NewEvent("length", l)
	}
	behavior := collector.New(10, processor)
	// Testing.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		tbe.Push(evt)
		switch evt.Topic() {
		case "length":
			var l int
			err := evt.Payload(&l)
			assert.NoError(err)
			assert.Equal(l, 10)
			tbe.SetSuccess()
		case behaviors.TopicResetted:
			return nil
		}
		return nil
	}
	// Run tests.
	tb := mesh.NewTestbed(behavior, test)
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
