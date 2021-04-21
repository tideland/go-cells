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

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"

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
	words := []string{"alpha", "beta", "", "d", "e"}
	eval := func(evt *mesh.Event) ([]string, error) {

	}
}

// EOF
