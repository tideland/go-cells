// Tideland Go Cells - Behaviors - Pairer - Unit Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package pairer_test // import "tideland.dev/go/cells/behaviors/pairer"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"

	"tideland.dev/go/cells/behaviors/pairer"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestSuccess verifies the successful finding of at least two matching.
func TestSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	pairerFunc := func(fstEvt, sndEvt *mesh.Event) bool {
		// Simply try to find the second event with the same
		// topic as the first.
		if fstEvt == nil {
			return true
		}
		return fstEvt.Topic() == sndEvt.Topic()
	}
	behavior := pairer.New(pairerFunc, time.Second)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			tbe.WaitFor(func() bool { return tbe.Len() > 0 })
			tbe.Do(func(i int, evt *mesh.Event) error {
				tbe.Assert(evt.Topic() == pairer.TopicPairerMatch, "invalid topic %d: %v", i, evt)
				var pair pairer.Pair
				err := evt.Payload(&pair)
				tbe.Assert(err == nil, "error retrieving the pair event payload: %v", err)
				tbe.Assert(pair.First.Topic() == pair.Second.Topic(), "pair event topic not equal: %v", pair)
				return nil
			})
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 10000; i++ {
			topic := generator.LimitedWord(4, 5)
			out.Emit(topic)
		}
	}, 5*time.Second)
	assert.NoError(err)
}

// TestFailOneHit verifies the failing of finding a pair after already
// one has been found.
func TestFailOneHit(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	pairerFunc := func(fstEvt, sndEvt *mesh.Event) bool {
		// Only find first one.
		if fstEvt == nil {
			return true
		}
		return false
	}
	behavior := pairer.New(pairerFunc, 10*time.Millisecond)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			tbe.AssertRetry(func() bool { return tbe.Len() > 0 }, "no emitted events")
			tbe.Do(func(i int, evt *mesh.Event) error {
				tbe.Assert(evt.Topic() != pairer.TopicPairerMatch, "invalid topic %d: %v", i, evt)
				return nil
			})
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 1000; i++ {
			topic := generator.LimitedWord(4, 5)
			out.Emit(topic)
		}
	}, time.Second)
	assert.NoError(err)
}

// TestFailNoHit verifies the failing of finding any pair.
func TestFailNoHit(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	pairerFunc := func(fstEvt, sndEvt *mesh.Event) bool {
		// Never find any one.
		return false
	}
	behavior := pairer.New(pairerFunc, 10*time.Millisecond)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			tbe.Assert(tbe.Len() == 0, "no pairer output expected: %v", tbe)
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 10000; i++ {
			topic := generator.LimitedWord(4, 5)
			out.Emit(topic)
		}
	}, time.Second)
	assert.NoError(err)
}

// EOF
