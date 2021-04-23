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
	"fmt"
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

// TestCounterBehavior tests counting and reacting via the
// countung behavior.
func TestCounterBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	topics := []string{"alpha", "bravo", "charly", "delta", "echo"}
	counteval := func(evt *mesh.Event) ([]string, error) {
		var incrs []string
		for _, r := range evt.Topic() {
			reg := fmt.Sprintf("reg-%v", r)
			incrs = append(incrs, reg)
		}
		return incrs, nil
	}
	behavior := behaviors.NewCounterBehavior(counteval)
	// Test evaluation.
	eval := func(evt *mesh.Event) (bool, error) {
		switch evt.Topic() {
		case behaviors.TopicCounterValues:
			var values map[string]int
			err := evt.Payload(&values)
			if err != nil {
				return false, err
			}
			l := len(values)
			return l == 13 || l == 0, nil
		default:
			return false, nil
		}
	}
	// Run test.
	tb := mesh.NewTestbed(behavior, eval)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 50; i++ {
			topic := generator.OneStringOf(topics...)
			out.Emit(topic)
		}
		// Retrieve, erase, and retrieve again.
		out.Emit(behaviors.TopicCounterStatus)
		out.Emit(behaviors.TopicCounterReset)
		out.Emit(behaviors.TopicCounterStatus)
	}, time.Second)
	assert.NoError(err)
}

// EOF
