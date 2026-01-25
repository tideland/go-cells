// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package mesh

import (
	"context"
	"sync"
	"testing"

	"tideland.dev/go/asserts/verify"
)

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
		verify.NoError(t, err)
	}

	wg.Wait()
	cancel()
}
