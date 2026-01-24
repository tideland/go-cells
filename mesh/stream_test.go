// Tideland Go Cells - Mesh - Tests
//
// Copyright (C) 2010-2026 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh // import "tideland.dev/go/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"sync"
	"testing"

	"tideland.dev/go/asserts/verify"
)

//--------------------
// TESTS
//--------------------

// TestStreamSimple verifies simple emitting and pulling of events
// via a stream.
func TestStreamSimple(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	str := newStream()
	topics := []string{"one", "two", "three", "four", "five"}

	var wg sync.WaitGroup

	wg.Add(20)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case evt := <-str.Pull():
				verify.Contains(t, evt.Topic(), topics)
				wg.Done()
			}
		}
	}()

	for i := 0; i < 20; i++ {
		topic := topics[i%len(topics)]
		err := str.Emit(topic)
		verify.NoError(t,err)
	}

	wg.Wait()
	cancel()
}

// EOF
