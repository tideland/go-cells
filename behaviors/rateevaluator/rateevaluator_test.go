// Tideland Go Cells - Behaviors - Rate Evaluator - Unit Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package rateevaluator_test // import "tideland.dev/go/cells/behaviors/rateevaluator"

//--------------------
// IMPORTS
//--------------------

import (
	"strings"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"

	"tideland.dev/go/cells/behaviors/rateevaluator"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestSuccess verifies the successful finding of at least two matching.
func TestSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	raterFunc := func(evt *mesh.Event) (bool, error) {
		// Each topic starting with an 'a' fires the rater.
		return strings.IndexRune(evt.Topic(), 'a') == 0
	}
	behavior := rateevaluator.New(raterFunc)
	// Testing.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		return nil
	}
	// Run test.
	tb := mesh.NewTestbed(behavior, test)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 10000; i++ {
			topic := generator.LimitedWord(3, 6)
			out.Emit(topic)
		}
	}, 5*time.Second)
	assert.NoError(err)
}

// EOF
